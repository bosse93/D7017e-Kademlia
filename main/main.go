package main

import (
	"fmt"
	"strconv"
)

func main() {
	/*
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(contact)

	for i := 1; i < 100; i++ {
		port := 8000 + i
		s := []string{}
		s = append(s, "localhost:")
		s = append(s, strconv.Itoa(port))
		a := strings.Join(s, "")
		rt.AddContact(NewContact(NewKademliaID(NewRandomKademliaID().String()), a))
	}
	*/

	kademlia := NewKademlia()

	IDRTList := map[KademliaID]RoutingTable{}

	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	//create 100 nodes
	for i := 1; i < 100; i++ {
		port := 8000 + i
		a := "localhost" + strconv.Itoa(port)
		ID := NewKademliaID(NewRandomKademliaID().String())
		rt := NewRoutingTable(NewContact(ID, a))
		IDRTList[*ID] = *rt
	}

	//each node joins by doing a lookup on the first node and populating its own table
	for k, v := range IDRTList {

		v.AddContact(firstNode)
		closestNodes := kademlia.LookupContact(IDRTList[k].me, *firstNodeRT)
		for i := range closestNodes {
			v.AddContact(closestNodes[i])
		}
	}

	//print the table of the first node
	for i := range firstNodeRT.buckets {
		contactList := firstNodeRT.buckets[i]
		fmt.Println("Bucket: " + strconv.Itoa(i))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
}

/*
	c := make(chan []Contact)
	kademlia := NewKademlia(rt)
	go kademlia.LookupContact(contact, c)
	contacts := <-c
	//contacts := rt.FindClosestContacts(sample.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())

	}
*/
}