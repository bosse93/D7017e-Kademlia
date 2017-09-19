package main

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
	kademlia.rt = *rt
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target Contact, network map[KademliaID]*RoutingTable) {
	//channel for data returned to this func
	c := make(chan int)
	//channels that returns data to each thread
	kademlia.threadChannels[0] = make(chan []Contact)
	kademlia.threadChannels[1] = make(chan []Contact)
	kademlia.threadChannels[2] = make(chan []Contact)

	kademlia.closest = NewContactCandidates()
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target.ID, 20))
	//calls alpha lookuphelpers
	for i := 0; i < 3 && i < len(kademlia.closest.contacts); i++ {
		kademlia.LookupHelper(target, network, c, i, 0)
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
				for i := range kademlia.threadChannels {
					//maybe send empty contact-array or something to the channels and check for that to stop the recursion,
					//not sure if closing the channel is enough
					close(kademlia.threadChannels[i])
				}
			}
		}
	}
}

func (kademlia *Kademlia) LookupHelper(target Contact, network map[KademliaID]*RoutingTable, c chan int, thread int, round int)  {
	threadChannel := kademlia.threadChannels[thread]
	//start new thread
	for i := 0; i < 20; i++{
		if _, ok := kademlia.asked[*kademlia.closest.contacts[i].ID]; !ok {
			go network[*kademlia.closest.contacts[i].ID].FindClosestContactsChannel(target.ID, 20, threadChannel)
			kademlia.asked[*kademlia.closest.contacts[i].ID] = true
			break
		}
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
