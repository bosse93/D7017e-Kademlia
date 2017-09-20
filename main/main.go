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



	IDRTList := map[KademliaID]*RoutingTable{}

	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	//firstNode := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	IDRTList[*firstNode.ID] = firstNodeRT

	//kademlia := NewKademlia(firstNodeRT)

	/*nodeIDs := []string{"0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"F0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FFF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0FF",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0F",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0"}*/
	//create 100 nodes
	for i := 0; i < 50; i++ {
		port := 8001 + i
		a := "localhost" + strconv.Itoa(port)
		ID := NewRandomKademliaID()
		//ID := NewKademliaID(nodeIDs[i])
		rt := NewRoutingTable(NewContact(ID, a))
		IDRTList[*ID] = rt
	}

	//each node joins by doing a lookup on the first node and populating its own table
	h := 1
	for k, v := range IDRTList {
		if k != *firstNode.ID {

			fmt.Println("Ny Nod varv " + strconv.Itoa(h) + ": " + v.me.String())
			kademlia := NewKademlia(v)
			v.AddContact(firstNodeRT.me)

			lookupResult := kademlia.LookupContact(IDRTList[k].me.ID, IDRTList)
			//fmt.Println(lookupResult)
			for q := range lookupResult {
				v.AddContact(lookupResult[q])
			}
			//firstNodeRT.AddContact(IDRTList[k].me)
		}
		h++	
	}

/*
	//print the table of the first node
	fmt.Println("Node: " + firstNode.String())
	for i := range firstNodeRT.buckets {
		contactList := firstNodeRT.buckets[i]
		fmt.Println("Bucket: " + strconv.Itoa(i))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
	*/

	//print the table of all nodes
	
	for q, w := range IDRTList {
		fmt.Println("Node: " + q.String())
		for z := range w.buckets {
			contactList := w.buckets[z]
			fmt.Println("Bucket: " + strconv.Itoa(z))
			for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
				contact := elt.Value.(Contact)
				fmt.Println(contact.String())
			}
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