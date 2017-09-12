package main

type Kademlia struct {
}

func NewKademlia() *Kademlia {
	return &Kademlia{}
}

func (kademlia *Kademlia) LookupContact(target Contact, rt RoutingTable) []Contact{
	rt.AddContact(target)
	return rt.FindClosestContacts(target.ID, 20)

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
