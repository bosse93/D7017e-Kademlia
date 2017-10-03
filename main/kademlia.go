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
	rt *RoutingTable
	network Network
	numberOfIdenticalAnswersInRow int
	threadCount int
}

func NewKademlia(nw *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.asked = make(map[KademliaID]bool)
	kademlia.network = *nw
	kademlia.rt = kademlia.network.node.rt
	kademlia.numberOfIdenticalAnswersInRow = 0
	kademlia.threadCount = 0
	rand.Seed(time.Now().UnixNano())
	return kademlia
}

func (kademlia *Kademlia) findNextNodeToAsk() (nextContact *Contact, success bool) {
	for i := range kademlia.closest.contacts {
		if(kademlia.asked[*kademlia.closest.contacts[i].ID] != true) {
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

func (kademlia *Kademlia) AskNextNode(target *KademliaID, destination *Contact, findData bool, returnChannel chan interface{}) {
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
							kademlia.AskNextNode(target, destination, findData, returnChannel)
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
							kademlia.AskNextNode(target, destination, findData, returnChannel)
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
						kademlia.AskNextNode(target, destination, findData, returnChannel)
					} 
				}			
		}
	}
}

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
