package main

import (
	"fmt"
	"strconv"
	"strings"
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

	for i := 1; i < 100; i++ {
		port := 8000 + i
		s := []string{}
		s = append(s, "localhost:")
		s = append(s, strconv.Itoa(port))
		a := strings.Join(s, "")
		ID := NewKademliaID(NewRandomKademliaID().String())
		rt := NewRoutingTable(NewContact(ID, a))
		IDRTList[*ID] = *rt
	}

	/*
	for k, v := range IDRTList {
		fmt.Println("key[%s] value[%s]", k.String(),v.me.String())
	}
*/

	/*
	
	*/

	for k, v := range IDRTList {

		v.AddContact(firstNode)
		closestNodes := kademlia.LookupContact(IDRTList[k].me, *firstNodeRT)
		for i := range closestNodes {
			v.AddContact(closestNodes[i])
		}
	}

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