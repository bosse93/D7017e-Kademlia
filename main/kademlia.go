package main

type Kademlia struct {
	table *RoutingTable
}

func NewKademlia(table *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.table = table
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target Contact, c chan []Contact) {
	c <- kademlia.table.FindClosestContacts(target.ID, 20)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
