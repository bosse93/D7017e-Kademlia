package main

import (
	"fmt"
	"strconv"
)

type Kademlia struct {
	closest ContactCandidates
	asked map[KademliaID]bool
	round map[Round][]Contact
	threadChannels [3]chan []Contact
	rt RoutingTable
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
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target Contact, network map[KademliaID]*RoutingTable, r chan []Contact) {
	//channel for data returned to this func
	c := make(chan int)
	//channels that returns data to each thread
	fmt.Println("LookupContact")
	kademlia.closest = NewContactCandidates()

	var threads = 0

	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target.ID, 20)) //3 räcker?
	//calls alpha lookuphelpers
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		go kademlia.LookupHelper(target, network, c, i, 0)
		kademlia.threadChannels[i] = make(chan []Contact)
		threads++
	}
	//after one thread is done with one round, if all threads are done for that round compare with previous round.
	//if everyone returned the same as the previous round close the channels
	for {
		select {
		case round := <-c:
			fmt.Println("round" + strconv.Itoa(round))
			if round > 0 {
				_, t0 := kademlia.round[Round{round, 0}]
				_, t1 := kademlia.round[Round{round, 1}]
				_, t2 := kademlia.round[Round{round, 2}]
				var same = true
				if t0 && t1 && t2 {
					for i := 0; i < 3 && same; i++ {
						for j := 0; j < 20 && same; j++ {
							if kademlia.round[Round{round, i}][j] != kademlia.round[Round{round-1, i}][j] {
								same = false
							}
						}
					}
				}
				if same {
					count := 20
					if count > kademlia.closest.Len() {
						count = kademlia.closest.Len()
					}
					r <- kademlia.closest.GetContacts(count)
					close(c)
					for i := 0; i < threads; i++ {
						//maybe send empty contact-array or something to the channels and check for that to stop the recursion,
						//not sure if closing the channel is enough
						close(kademlia.threadChannels[i])
					}

				}
			} else if round == -1 {
				fmt.Println("return")
				r <- kademlia.closest.GetContacts(20)
				close(c)
				for i := 0; i < threads; i++ {
					//maybe send empty contact-array or something to the channels and check for that to stop the recursion,
					//not sure if closing the channel is enough
					/*var contact []Contact
					contact[0] = NewContact(NewKademliaID("quit"), "quit")
					kademlia.threadChannels[i] <- contact*/
					close(kademlia.threadChannels[i])
				}


			}
		}
	}
}

func (kademlia *Kademlia) LookupHelper(target Contact, network map[KademliaID]*RoutingTable, c chan int, thread int, round int)  {
	threadChannel := kademlia.threadChannels[thread]
	//start new thread
	for i := range kademlia.closest.contacts {
		fmt.Println("closest" + kademlia.closest.contacts[i].ID.String() + " round " + strconv.Itoa(round))
	}

	for i := 0; i < 20 && i < len(kademlia.closest.contacts); i++{
		if _, ok := kademlia.asked[*kademlia.closest.contacts[i].ID]; !ok {
			table := network[*kademlia.closest.contacts[i].ID]
			go table.FindClosestContactsChannel(target.ID, 20, threadChannel)
			kademlia.asked[*kademlia.closest.contacts[i].ID] = true
			break
		}
		if i == len(kademlia.closest.contacts) - 1 {
			fmt.Println("done")
			c <- -1
		}
		//Om i har itererat igenom alla contacter i closest
		//contacts utan att hittat nån som inte blivit tillfrågad ännu
		//Vad göra? Invänta alla andra trådar? Avsluta funktionen och därmed rekursionen?
	}
	//update info, notify channel and do recursive call
	select {
	case x := <-threadChannel:
		//if (x[0].Address != "quit") {
			kademlia.closest.Append(x)
			kademlia.closest.Sort()
			kademlia.round[Round{round, thread}] = x
			c <- round
			kademlia.LookupHelper(target, network, c, thread, round+1)
		//}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
