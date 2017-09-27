package main

import (
	"time"
	//"fmt"
	"math/rand"
	//"strconv"
	"fmt"
	"strconv"
)

type Kademlia struct {
	closest ContactCandidates
	asked map[KademliaID]bool
	threadChannels [3]chan []Contact
	returnDataChannels map[KademliaID](chan string)
	rt *RoutingTable
	network Network
	numberOfIdenticalAnswersInRow int
	noMoreNodesTimeout int
	done bool
	threadCount int
}


func NewKademlia(nw *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.asked = make(map[KademliaID]bool)
	kademlia.returnDataChannels = make(map[KademliaID]chan string)
	kademlia.network = *nw
	kademlia.rt = kademlia.network.node.rt
	kademlia.numberOfIdenticalAnswersInRow = 0
	kademlia.noMoreNodesTimeout = 0
	kademlia.done = false
	kademlia.threadCount = 0
	rand.Seed(time.Now().UnixNano())
	return kademlia
}


func (kademlia *Kademlia) findNextNodeToAsk() (nextContact *Contact, success bool) {
	for i := range kademlia.closest.contacts {
		//fmt.Println("NextNode loop: " + kademlia.closest.contacts[i].ID.String())
		if(kademlia.asked[*kademlia.closest.contacts[i].ID] != true) {
			fmt.Println("Found new Node")
			kademlia.asked[*kademlia.closest.contacts[i].ID] = true
			nextContact = &kademlia.closest.contacts[i]
			success = true
			return 
		}
	}
	nextContact = nil
	success = false
	return
}

func (kademlia *Kademlia) askNextNode(target *KademliaID, destination *Contact, findData bool, returnChannel chan interface{}) {
		fmt.Println("AskNextNode Success")
		if(findData) {
			go kademlia.network.SendFindDataMessage(target.String(), destination, returnChannel)
		} else {
			go kademlia.network.SendFindContactMessage(target, destination, returnChannel)
		}
}

func (kademlia *Kademlia) updateClosestContacts(networkAnswer []Contact, target *KademliaID) {
	same := true
	var newNodeList []Contact
	for i := range networkAnswer {
		existsAlready := false
		for k := range kademlia.closest.contacts {
			if(networkAnswer[i].ID.String() == kademlia.closest.contacts[k].ID.String()) {	
				existsAlready = true
			}
		}
		if(!existsAlready) {
			same = false
			networkAnswer[i].CalcDistance(target)
			newNodeList = append(newNodeList, networkAnswer[i])
		}
	}

	if(same) {
		kademlia.numberOfIdenticalAnswersInRow++
	} else {
		kademlia.numberOfIdenticalAnswersInRow = 0
	}

	kademlia.closest.Append(newNodeList)
	kademlia.closest.Sort()

	numberOfResults := 20
	if (len(kademlia.closest.contacts) < 20) {
		numberOfResults = len(kademlia.closest.contacts)
	}
	kademlia.closest.contacts = kademlia.closest.GetContacts(numberOfResults)
}

func (kademlia *Kademlia) LookupContact(target *KademliaID, findData bool) (returnContact []Contact, dataReturn string) {
	kademlia.closest = NewContactCandidates()
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target, 3))

	returnChannel := make(chan interface{}, 3)
	
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		if findData {
			fmt.Println("New Thread")
			kademlia.threadCount++
			go kademlia.network.SendFindDataMessage(target.String(), &kademlia.closest.contacts[i], returnChannel)
		} else {
			fmt.Println("New Thread")
			kademlia.threadCount++
			go kademlia.network.SendFindContactMessage(target, &kademlia.closest.contacts[i], returnChannel)
		}
		kademlia.asked[*kademlia.closest.contacts[i].ID] = true
	}
	for {
		select {
			case networkAnswer := <-returnChannel:
				switch networkAnswer := networkAnswer.(type) {
					case []Contact:
						kademlia.updateClosestContacts(networkAnswer, target)
						if(kademlia.numberOfIdenticalAnswersInRow > 2) {
							returnContact = kademlia.closest.contacts
							fmt.Println("Same Answer in a row")
							return 
						}
						destination, success := kademlia.findNextNodeToAsk()
						if(success) {
							kademlia.askNextNode(target, destination, findData, returnChannel)
						} else {
							fmt.Println("Thread Killed")
							kademlia.threadCount--
						}

					case string:
						returnContact = []Contact{}
						dataReturn = networkAnswer
						return  

					case bool:
						fmt.Println("Timeout")
						destination, success := kademlia.findNextNodeToAsk()
						if(success) {
							kademlia.askNextNode(target, destination, findData, returnChannel)
						} else {
							fmt.Println("Thread Killed")
							kademlia.threadCount--
						}
				}
	
			default:
				if(kademlia.threadCount == 0) {
					fmt.Println("No Threads")
					returnContact = kademlia.closest.contacts
					return
				}
				if(kademlia.threadCount < 3) {
					destination, success := kademlia.findNextNodeToAsk()
					if(success) {
						fmt.Println("New Thread")
						kademlia.threadCount++
						kademlia.askNextNode(target, destination, findData, returnChannel)
					} 
				}			
		}
	}
}


/*
func (kademlia *Kademlia) LookupContact(target *KademliaID, findData bool) []Contact {
	kademlia.closest = NewContactCandidates()
	//Find up to alpha contacts closest to target in own RT
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target, 3)) //3 räcker?


	//Start buffered channel (Async) for communication mainThread -> askThreads
	destinationChannel := make(chan Contact, 6)
	
	//Start up to alpha threads on unique buffered channels (Async).
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		kademlia.threadChannels[i] = make(chan []Contact, 2)
		go kademlia.LookupHelper(target, kademlia.closest.contacts[i], kademlia.threadChannels[i], destinationChannel, kademlia.threadCount + 1, findData)
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
						go kademlia.LookupHelper(target, destinationContact, kademlia.threadChannels[kademlia.threadCount], destinationChannel, kademlia.threadCount + 1, findData)
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


	

func (kademlia *Kademlia) LookupHelper(target *KademliaID, destination Contact, sendChannel chan []Contact, recieveChannel chan Contact, threadNmbr int, findData bool)  {
	findContactReturn := []Contact{}
	if findData {
		dataReturn := kademlia.network.SendFindDataMessage(target.String(), &destination)
		fmt.Println("first id returned: " + dataReturn.contacts[0].ID.String() + " Data: " + dataReturn.data)
		findContactReturn = dataReturn.contacts
		if dataReturn.data != "" {
			fmt.Println("got data in lookupHelper: " + dataReturn.data)
			kademlia.returnDataChannels[*target] <- dataReturn.data
		}

	} else {
		findContactReturn = kademlia.network.SendFindContactMessage(&destination, target)
	}

	for i := range findContactReturn {
		findContactReturn[i].CalcDistance(target)
	}
	sendChannel <-findContactReturn
	select {
	//If channel is closed main thread have decided Lookup is done. Close own channel and end recursion.
	case nextDestination, ok := <-recieveChannel:
		if ok {
			if ((findContactReturn[0].ID.String() != "0000000000000000000000000000000000000000") && (findContactReturn[0].Address != "0.0.0.0:0000")) {
			//Found contact in destination channel. Ask the node!
				kademlia.LookupHelper(target, nextDestination, sendChannel, recieveChannel, threadNmbr, findData)
			}
		} else {
			close(sendChannel)
			break
		}
	}
}

func (kademlia *Kademlia) answerHelper(answer []Contact) {
	if(len(answer) == 0 || ((answer[0].ID.String() == "0000000000000000000000000000000000000000") && (answer[0].Address == "0.0.0.0:0000"))){
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
*/

func (kademlia *Kademlia) LookupData(hash string) (returnData string, success bool) {
	contacts, data := kademlia.LookupContact(NewKademliaID(hash), true)
	if(len(contacts) == 0) {
		fmt.Println("LookupData found data")
		returnData = data
		success = true
		return
	} else {
		success = false
		fmt.Println("LookupData did not find data")
		return
	}
}

func (kademlia *Kademlia) Store(key *KademliaID, data string) {
	contacts, _ := kademlia.LookupContact(key, false)
	for i := 0 ; i < len(contacts); i++ {
		returnChannel := make(chan interface{})
		go kademlia.network.SendStoreMessage(key.String(), data, contacts[i].Address, returnChannel)
		returnValue:= <-returnChannel
		switch returnValue := returnValue.(type) {
			case string:
				fmt.Println("Store " + strconv.Itoa(i) + " Reply: " + returnValue)
			default:
				fmt.Println("Something went wrong")
		}
	}
}
