package main

import (
	"fmt"
	"strconv"
	"strings"
	sample "D7017e-Kademlia/SampleCode"
)

func main() {
	rt := sample.NewRoutingTable(sample.NewContact(sample.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	for i := 1; i < 100; i++ {
		port := 8000 + i
		s := []string{}
		s = append(s, "localhost:")
		s = append(s, strconv.Itoa(port))
		a := strings.Join(s, "")
		rt.AddContact(sample.NewContact(sample.NewKademliaID(sample.NewRandomKademliaID().String()), a))
	}

	contacts := rt.FindClosestContacts(sample.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}