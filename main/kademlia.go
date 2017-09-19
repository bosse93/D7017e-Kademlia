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
	kademlia.closest = NewContactCandidates()

	var threads = 0

	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target.ID, 20)) //3 räcker?
	/*
	for i := range kademlia.rt.buckets {
		contactList := kademlia.rt.buckets[i]
		fmt.Println("Bucket: " + strconv.Itoa(i))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
	*/
	
	//calls alpha lookuphelpers
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {

		//kademlia.LookupHelper(target, network, c, i, 0)

		go kademlia.LookupHelper(target, network, c, i, 0)
		kademlia.threadChannels[i] = make(chan []Contact)
		threads++
	}
	//after one thread is done with one round, if all threads are done for that round compare with previous round.
	//if everyone returned the same as the previous round close the channels
	select {
	case round := <-c:
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
				close(c)
				for i := 0; i < threads; i++ {
					//maybe send empty contact-array or something to the channels and check for that to stop the recursion,
					//not sure if closing the channel is enough
					close(kademlia.threadChannels[i])
				}
				r <- kademlia.closest.GetContacts(20)
			}
		}
	}
}

func (kademlia *Kademlia) LookupHelper(target Contact, network map[KademliaID]*RoutingTable, c chan int, thread int, round int)  {
	threadChannel := kademlia.threadChannels[thread]
	//start new thread
	for i := 0; i < 20; i++{
		if i < len(kademlia.closest.contacts) {
			if _, ok := kademlia.asked[*kademlia.closest.contacts[i].ID]; !ok {
				go network[*kademlia.closest.contacts[i].ID].FindClosestContactsChannel(target.ID, 20, threadChannel)
				kademlia.asked[*kademlia.closest.contacts[i].ID] = true
				break
			}
		//Om i har itererat igenom alla contacter i closest
		//contacts utan att hittat nån som inte blivit tillfrågad ännu
		//Vad göra? Invänta alla andra trådar? Avsluta funktionen och därmed rekursionen?
	}
	//update info, notify channel and do recursive call
	select {
	case x := <-threadChannel:
		kademlia.closest.Append(x)
		kademlia.closest.Sort()
		kademlia.round[Round{round, thread}] = x
		c <- round
		kademlia.LookupHelper(target, network, c, thread, round+1)
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
