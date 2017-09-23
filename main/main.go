package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	nw := NewNetwork(NewNode(firstNodeRT))
	go nw.Listen("localhost", 8000)

	time.Sleep(500 * time.Millisecond)	
	nodeList := []*RoutingTable{firstNodeRT}
	
	//create 100 nodes
	for i := 0; i < 100; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)

		
		ID := NewRandomKademliaID()
		rt := NewRoutingTable(NewContact(ID, a))
		nodeList = append(nodeList, rt)
		rt.AddContact(firstNodeRT.me)
		nw := NewNetwork(rt)
		fmt.Println("Ny Nod varv " + strconv.Itoa(i+1) + ": " + rt.me.String())
		go nw.Listen("localhost", port)
		time.Sleep(500 * time.Millisecond)
		kademlia := NewKademlia(nw)

		lookupResult := kademlia.LookupContact(ID)

		for q := range lookupResult {
				rt.AddContact(lookupResult[q])
		}	
			
	}

	printFirstNodeRT(firstNode, firstNodeRT)
	printLastNodeRT(nodeList)
	
}

func printFirstNodeRT(firstNode Contact, firstNodeRT *RoutingTable) {
	fmt.Println("Node: " + firstNode.ID.String())
	for z := range firstNodeRT.buckets {
		contactList := firstNodeRT.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
}

func printLastNodeRT(nodeList []*RoutingTable) {
	lastNode := nodeList[len(nodeList)-1]
	fmt.Println("Node: " + lastNode.me.String())
	for z := range lastNode.buckets {
		contactList := lastNode.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
}

func printAllNodesRT(nodeList []*RoutingTable) {
	for w := range nodeList {
		fmt.Println("Node: " + nodeList[w].me.String())
		for z := range nodeList[w].buckets {
			contactList := nodeList[w].buckets[z]
			fmt.Println("Bucket: " + strconv.Itoa(z))
			for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
				contact := elt.Value.(Contact)
				fmt.Println(contact.String())
			}
		}
	}
}
