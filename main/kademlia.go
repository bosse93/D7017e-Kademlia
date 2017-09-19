package main

type Kademlia struct {
	closest ContactCandidates
	asked map[KademliaID]bool
	rt RoutingTable
}

func NewKademlia(rt *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = *rt
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target Contact, network map[KademliaID]*RoutingTable) {
	c := make(chan []Contact)
	kademlia.closest = NewContactCandidates()
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target.ID, 20))
	for i := 0; i < 3; i++ {
		if i < len(kademlia.closest.contacts) {
			kademlia.LookupHelper(target, network, c)
		}
	}
}

func (kademlia *Kademlia) LookupHelper(target Contact, network map[KademliaID]*RoutingTable, c chan []Contact)  {
	for i := 0; i < 20; i++{
		if _, ok := kademlia.asked[*kademlia.closest.contacts[i].ID]; !ok {
			go network[*kademlia.closest.contacts[i].ID].FindClosestContactsChannel(target.ID, 20, c)
			kademlia.asked[*kademlia.closest.contacts[i].ID] = true
			/*
			if true {
				break
			}
			*/
		}
	}
	select {
	case x := <-c:
		kademlia.closest.Append(x)
		kademlia.closest.Sort()
		kademlia.LookupHelper(target, network, c)
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
