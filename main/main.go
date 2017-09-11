package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
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
	c := make(chan []Contact)
	kademlia := NewKademlia(rt)
	go kademlia.LookupContact(contact, c)
	contacts := <-c
	//contacts := rt.FindClosestContacts(sample.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}