package main

import (
	"time"
	//"fmt"
	"math/rand"
	//"strconv"
)

type Kademlia struct {
	closest ContactCandidates
	asked map[KademliaID]bool
	round map[Round][]Contact
	threadChannels [3]chan []Contact
	rt *RoutingTable
	networkTest Network
	numberOfIdenticalAnswersInRow int
	noMoreNodesTimeout int
	done bool
	threadCount int
}

type Round struct {
	round int
	thread int
}

func NewKademlia(nw *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.asked = make(map[KademliaID]bool)
	kademlia.round = make(map[Round][]Contact)
	kademlia.networkTest = *nw
	kademlia.rt = kademlia.networkTest.node.rt
	kademlia.numberOfIdenticalAnswersInRow = 0
	kademlia.noMoreNodesTimeout = 0
	kademlia.done = false
	kademlia.threadCount = 0
	rand.Seed(time.Now().UnixNano())
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	kademlia.closest = NewContactCandidates()
	//Find up to alpha contacts closest to target in own RT
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target, 3)) //3 räcker?


	//Start buffered channel (Async) for communication mainThread -> askThreads
	destinationChannel := make(chan Contact, 6)
	
	//Start up to alpha threads on unique buffered channels (Async).
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		kademlia.threadChannels[i] = make(chan []Contact, 2)
		go kademlia.LookupHelper(target, kademlia.closest.contacts[i], kademlia.threadChannels[i], destinationChannel, kademlia.threadCount + 1)
		kademlia.threadCount++
		kademlia.asked[*kademlia.closest.contacts[i].ID] = true
	}

	//Loop until Lookup is done.
	for {
		//Check if anyone of c1, c2, c3 is ready. If non is ready do default.
		select {
			case c1 := <-kademlia.threadChannels[0]:
				//fmt.Println("Channel 1")
				kademlia.answerHelper(c1)

			case c2 := <-kademlia.threadChannels[1]:
				//fmt.Println("Channel 2")
				kademlia.answerHelper(c2)

			case c3 := <-kademlia.threadChannels[2]:
				//fmt.Println("Channel 3")
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
						go kademlia.LookupHelper(target, destinationContact, kademlia.threadChannels[kademlia.threadCount], destinationChannel, kademlia.threadCount + 1)
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


	

func (kademlia *Kademlia) LookupHelper(target *KademliaID, destination Contact, sendChannel chan []Contact, recieveChannel chan Contact, threadNmbr int)  {
	//Sleep random. Simulating a network call to destination.
	//sleepTime := rand.Intn(20)
	//time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	//Should be network call and add to channel what the destination node answers.
	//Channel buffer is 2. If channel is full will wait until mainThread have recieved at least one answer.
	//fmt.Println("Thread " + strconv.Itoa(threadNmbr) + ":" + " Sent ping to " + destination.Address) 
	findContactReturn := kademlia.networkTest.SendFindContactMessage(&destination, target)
	//fmt.Println("Thread " + strconv.Itoa(threadNmbr) + ": Recieved answer from " + pingReturn)
	//sendChannel <-network[*destination.ID].rt.FindClosestContacts(target, 20)
	for i := range findContactReturn {
		findContactReturn[i].CalcDistance(target)
	}
	sendChannel <-findContactReturn
	//Add asker node to the asked ones RT. Should not be done here but in the other nodes client/RT
	//network[*destination.ID].rt.AddContact(kademlia.rt.me)
	select {
		//If channel is closed main thread have decided Lookup is done. Close own channel and end recursion.
		case nextDestination, ok := <-recieveChannel:
			if ok {
				//Found contact in destination channel. Ask the node!
				kademlia.LookupHelper(target, nextDestination, sendChannel, recieveChannel, threadNmbr)
			} else {
				close(sendChannel)
				break
			}
	}
}

func (kademlia *Kademlia) answerHelper(answer []Contact) {
	if(len(answer) == 0) {
		return
	}
	same := true
	var newNodeList []Contact
	//Check if elements in answer already exists in closest. If so do not add it again.
	
	for i := range answer {
		existsAlready := false
		for k := range kademlia.closest.contacts {
			if(answer[i].ID.String() == kademlia.closest.contacts[k].ID.String()) {	
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

func (kademlia *Kademlia) Store(key *KademliaID, data string, network map[KademliaID]*Network) {
	contacts := kademlia.LookupContact(key)
	for i := 0 ; i < len(contacts); i++ {
		kademlia.networkTest.SendStoreMessage(key.String(), data, contacts[i].Address)
	}
}
