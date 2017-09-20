package main

import (
	"time"
	"fmt"
	"math/rand"
)

type Kademlia struct {
	closest ContactCandidates
	asked map[KademliaID]bool
	round map[Round][]Contact
	threadChannels [3]chan []Contact
	rt RoutingTable
	numberOfIdenticalAnswersInRow int
	noMoreNodesTimeout int
	done bool
	threadCount int
}

type Round struct {
	round int
	thread int
}

func NewKademlia(rt *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.asked = make(map[KademliaID]bool)
	kademlia.round = make(map[Round][]Contact)
	kademlia.rt = *rt
	kademlia.numberOfIdenticalAnswersInRow = 0
	kademlia.noMoreNodesTimeout = 0
	kademlia.done = false
	kademlia.threadCount = 0
	rand.Seed(time.Now().UnixNano())
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID, network map[KademliaID]*RoutingTable) []Contact {
	kademlia.closest = NewContactCandidates()
	//Find up to alpha contacts closest to target in own RT
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target, 3)) //3 räcker?


	//Start buffered channel (Async) for communication mainThread -> askThreads
	destinationChannel := make(chan Contact, 6)
	
	//Start up to alpha threads on unique buffered channels (Async).
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		kademlia.threadChannels[i] = make(chan []Contact, 2)
		go kademlia.LookupHelper(target, kademlia.closest.contacts[i], network, kademlia.threadChannels[i], destinationChannel)
		kademlia.threadCount++
		kademlia.asked[*kademlia.closest.contacts[i].ID] = true
	}

	//Loop until Lookup is done.
	for {
		//Check if anyone of c1, c2, c3 is ready. If non is ready do default.
		select {
			case c1 := <-kademlia.threadChannels[0]:
				fmt.Println("Channel 1")
				kademlia.answerHelper(c1)

			case c2 := <-kademlia.threadChannels[1]:
				fmt.Println("Channel 2")
				kademlia.answerHelper(c2)

			case c3 := <-kademlia.threadChannels[2]:
				fmt.Println("Channel 3")
				kademlia.answerHelper(c3)
			
			//Check if done cause of timeout or same answers in a row. If so end lookup.	
			default:
				if kademlia.done {
					break
				}
				if kademlia.numberOfIdenticalAnswersInRow > 2 {
					close(destinationChannel)
					kademlia.done = true
					numberOfResults := 20
					if len(kademlia.closest.contacts) < 20 {
						numberOfResults = len(kademlia.closest.contacts)
					}
					return kademlia.closest.GetContacts(numberOfResults)
				}
				//Check for node in closest which have not been asked already
				nodeFound := false
				destinationContact := NewContact(NewRandomKademliaID(), "None")
				for i := range kademlia.closest.contacts {
					if kademlia.asked[*kademlia.closest.contacts[i].ID] != true {
						destinationContact = kademlia.closest.contacts[i]
						nodeFound = true
						break
					}
				}
				//If node which not been asked already is found. Add it to channel for Main -> askThread communication
				if nodeFound {
					kademlia.noMoreNodesTimeout = 0
					if kademlia.threadCount < 3 {
						kademlia.threadChannels[kademlia.threadCount] = make(chan []Contact, 2)
						go kademlia.LookupHelper(target, destinationContact, network, kademlia.threadChannels[kademlia.threadCount], destinationChannel)
						kademlia.threadCount++
						kademlia.asked[*destinationContact.ID] = true
					} else {
						select {
							case destinationChannel <- destinationContact:
								kademlia.asked[*destinationContact.ID] = true
								break
							default:
								time.Sleep(5 * time.Millisecond)
						}
					}
				} else {
					//No node found. Timeout if no node found in some time.
					time.Sleep(5 * time.Millisecond)
					kademlia.noMoreNodesTimeout++
					if (kademlia.noMoreNodesTimeout > 10) {
						//fmt.Println("Timeout")
						close(destinationChannel)
						kademlia.done = true
						
						numberOfResults := 20
						if len(kademlia.closest.contacts) < 20 {
							numberOfResults = len(kademlia.closest.contacts)
						}
						return kademlia.closest.GetContacts(numberOfResults)
					}
				}		
		}
		
	}


}


	

func (kademlia *Kademlia) LookupHelper(target *KademliaID, destination Contact, network map[KademliaID]*RoutingTable, sendChannel chan []Contact, recieveChannel chan Contact)  {
	//Sleep random. Simulating a network call to destination.
	sleepTime := rand.Intn(20)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	//Should be network call and add to channel what the destination node answers.
	sendChannel <-network[*destination.ID].FindClosestContacts(target, 20)
	network[*destination.ID].AddContact(kademlia.rt.me)
	select {
		//If channel is closed main thread have decided Lookup is done. Close own channel and end recursion.
		case nextDestination, ok := <-recieveChannel:
			if ok {
				//Found contact in destination channel. Ask the node!
				kademlia.LookupHelper(target, nextDestination, network, sendChannel, recieveChannel)
			} else {
				close(sendChannel)
				break
			}
	}
}

func (kademlia *Kademlia) answerHelper(answer []Contact) {
	same := true
	var newNodeList []Contact
	//Check if elements in answer already exists in closest. If so do not add it again.
	for i := range answer {
		existsAlready := false
		for k := range kademlia.closest.contacts {
			if(answer[i].ID == kademlia.closest.contacts[k].ID) {
				existsAlready = true
			}
		}
		if(!existsAlready) {
			same = false
			newNodeList = append(newNodeList, answer[i])
		}
	}
	//Same answer as closest. If x identical answers in row lookup will decide its done
	if(same) {
		kademlia.numberOfIdenticalAnswersInRow++
	} else {
		kademlia.numberOfIdenticalAnswersInRow = 0
	}

	//Append new nodes to closest. Sort and limit it to the 20 best.
	kademlia.closest.Append(newNodeList)
	kademlia.closest.Sort()
	
	numberOfResults := 20
	if (len(kademlia.closest.contacts) < 20) {
		numberOfResults = len(kademlia.closest.contacts)
	}
	newCandidates := kademlia.closest.GetContacts(numberOfResults)
	kademlia.closest = NewContactCandidates()
	kademlia.closest.Append(newCandidates)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
