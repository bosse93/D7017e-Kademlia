package main

type Kademlia struct {
	closest ContactCandidates
	asked ContactCandidates
	rt RoutingTable
}

func NewKademlia(rt *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = *rt
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target Contact, network map[KademliaID]*RoutingTable) {
	kademlia.closest = NewContactCandidates()
	kademlia.asked = NewContactCandidates()
	kademlia.closest.Append(kademlia.rt.FindClosestContacts(target.ID, 20))
	for i := 0; i < 3; i++ {
		if i < len(kademlia.closest.contacts) {
			go network[*kademlia.closest.contacts[i].ID].FindClosestContacts(target.ID, 20)
		}
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
