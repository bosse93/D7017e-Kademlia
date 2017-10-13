package main

import (
	"time"
	//"fmt"
	"math/rand"
	//"strconv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Kademlia struct {
	// closest contains closest contacts found so far.
	closest ContactCandidates
	// asked contains true if id already been asked.
	asked                         map[KademliaID]bool
	rt                            *RoutingTable
	network                       Network
	numberOfIdenticalAnswersInRow int
	threadCount                   int
	// k number of network requests will be done simultaneously.
	k int
}

// NewKademlia initializes Kademlia object.
// Sets k value to the desired amount of simultaneously network requests.
func NewKademlia(nw *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.asked = make(map[KademliaID]bool)
	kademlia.network = *nw
	kademlia.rt = kademlia.network.node.rt
	kademlia.numberOfIdenticalAnswersInRow = 0
	kademlia.threadCount = 0
	kademlia.k = 3
	rand.Seed(time.Now().UnixNano())
	return kademlia
}

// FindNextNodeToAsk helps the Lookup function to find the next node to ask.
// This is done by finding a node which haven't been asked yet in the kademlia asked map.
func (kademlia *Kademlia) FindNextNodeToAsk() (nextContact *Contact, success bool) {
	for i := range kademlia.closest.contacts {
		if kademlia.asked[*kademlia.closest.contacts[i].ID] != true {
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

// AskNextNode will create a new go routine to do a network request.
// If findData is true a SendFindDataMessage will be done else SendFindContactMessage is done.
func (kademlia *Kademlia) AskNextNode(target *KademliaID, destination *Contact, findData bool, returnChannel chan interface{}) {
	if findData {
		go kademlia.network.SendFindDataMessage(target.String(), destination, returnChannel)
	} else {
		go kademlia.network.SendFindContactMessage(target, destination, returnChannel)
	}
}

// UpdateClosestContacts will update closest contacts with newly aquired contacts.
// All nodes that doesn't already exist in closest will be added.
// Then the list is sorted and upto 20 closest contacts is kept.
// If no new Contact is added to closest 1 is added to numberOfIdenticalAnswersInRow.
func (kademlia *Kademlia) UpdateClosestContacts(networkAnswer []Contact, target *KademliaID) {
	same := true
	var newNodeList []Contact
	for i := range networkAnswer {
		existsAlready := false
		for k := range kademlia.closest.contacts {
			if networkAnswer[i].ID.String() == kademlia.closest.contacts[k].ID.String() {
				existsAlready = true
			}
		}
		if !existsAlready {
			same = false
			networkAnswer[i].CalcDistance(target)
			newNodeList = append(newNodeList, networkAnswer[i])
		}
	}

	if same {
		kademlia.numberOfIdenticalAnswersInRow++
	} else {
		kademlia.numberOfIdenticalAnswersInRow = 0
	}

	kademlia.closest.Append(newNodeList)
	kademlia.closest.Sort()

	numberOfResults := 20
	if len(kademlia.closest.contacts) < 20 {
		numberOfResults = len(kademlia.closest.contacts)
	}
	kademlia.closest.contacts = kademlia.closest.GetContacts(numberOfResults)
}

// LookupContact returns 20 closest contacts to target.
// Spawns threads to do requests to Contacts and then handle their response.
// Stops if same answer is recieved multiple times or if all contacts in closest have been asked.
// If findData is set to true it will try to find data if it exists on the network.
func (kademlia *Kademlia) LookupContact(target *KademliaID, findData bool) (returnContact []Contact, dataReturn string) {
	kademlia.closest = NewContactCandidates()
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target, 3))

	returnChannel := make(chan interface{}, 3)

	for i := 0; i < 3 && i < len(kademlia.closest.contacts) && kademlia.threadCount < kademlia.k; i++ {
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
				kademlia.UpdateClosestContacts(networkAnswer, target)
				if kademlia.numberOfIdenticalAnswersInRow > 2 {
					returnContact = kademlia.closest.contacts
					fmt.Println("Same Answer in a row")
					return
				}
				destination, success := kademlia.FindNextNodeToAsk()
				if success {
					kademlia.AskNextNode(target, destination, findData, returnChannel)
				} else {
					fmt.Println("Thread Killed")
					kademlia.threadCount--
				}

			case string:
				fmt.Println(networkAnswer)
				returnContact = []Contact{}
				dataReturn = networkAnswer
				return

			case bool:
				fmt.Println("Timeout")
				destination, success := kademlia.FindNextNodeToAsk()
				if success {
					kademlia.AskNextNode(target, destination, findData, returnChannel)
				} else {
					fmt.Println("Thread Killed")
					kademlia.threadCount--
				}
			}

		default:
			if kademlia.threadCount == 0 {
				fmt.Println("No Threads")
				returnContact = kademlia.closest.contacts
				return
			}
			if kademlia.threadCount < kademlia.k {
				destination, success := kademlia.FindNextNodeToAsk()
				if success {
					fmt.Println("New Thread")
					kademlia.threadCount++
					kademlia.AskNextNode(target, destination, findData, returnChannel)
				}
			}
		}
	}
}

// LookupData will try to find data with name fileName in network.
// FileName hash is passed to LookupContact which will locate the data.
// When/if data is located fileNetwork.DownloadFile will download the file.
func (kademlia *Kademlia) LookupData(fileName string) bool {
	fileNameHash := HashKademliaID(fileName)

	//KIKA OM DATAN REDAN FINNS I STORAGE

	contacts, data := kademlia.LookupContact(fileNameHash, true)
	if len(contacts) == 0 {
		fmt.Println("LookupData found data")
		go kademlia.network.fileNetwork.DownloadFile(fileNameHash, data, true)
		return true
	} else {
		fmt.Println("LookupData did not find data")
		return false
	}
}

// Store will store fileName on network.
// Uses LookupContact to find closest contacts to hash of fileName.
// When closest contacts is found a store message is sent through SendStoreAndWaitForAnswer.
func (kademlia *Kademlia) Store(fileName string) {
	fileNameHash := HashKademliaID(fileName)
	contacts, _ := kademlia.LookupContact(fileNameHash, false)
	for i := 0; i < len(contacts); i++ {
		if contacts[i].ID.String() != kademlia.rt.me.ID.String() {
			go kademlia.SendStoreAndWaitForAnswer(fileNameHash.String(), contacts[i].Address, i)
		} else {
			fileDst, _ := os.Create("kademliastorage/" + kademlia.rt.me.ID.String() + "/" + fileNameHash.String())
			fileSrc, _ := os.Open("upload/" + kademlia.rt.me.ID.String() + "/" + fileName)

			if _, err := io.Copy(fileDst, fileSrc); err != nil {
				log.Fatal(err)
			}
			kademlia.network.node.Store(*fileNameHash, time.Now())
		}
	}
}

// SendStoreAndWaitForAnswer sends a store message to address and await answer.
// If store request is successfully accepted a confirmation will be recieved.
func (kademlia *Kademlia) SendStoreAndWaitForAnswer(fileName string, address string, number int) {
	returnChannel := make(chan interface{})
	go kademlia.network.SendStoreMessage(fileName, address, returnChannel)
	returnValue := <-returnChannel
	switch returnValue := returnValue.(type) {
	case string:
		fmt.Println("Store " + strconv.Itoa(number) + " Reply: " + returnValue)
	case bool:
		fmt.Println("Store request timeout")
	default:
		fmt.Println("Something went wrong")
	}
}
