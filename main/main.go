package main

import (
	"fmt"
	"strconv"
	//"log"
	//"D7024e-Kademlia/protobuf/proto"
	"time"
)

func main() {
	IDRTList := map[KademliaID]*Network{}

	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	//firstNode := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	nw := NewNetwork(NewNode(firstNodeRT))
	go nw.Listen("localhost", 8000)
	IDRTList[*firstNode.ID] = nw

	//kademlia := NewKademlia(firstNodeRT)

	//create 100 nodes
	for i := 0; i < 100; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)
		ID := NewRandomKademliaID()
		//ID := NewKademliaID(nodeIDs[i])

		rt := NewRoutingTable(NewContact(ID, a))
		nw := NewNetwork(NewNode(rt))
		go nw.Listen("localhost", port)
		IDRTList[*ID] = nw
	}
	time.Sleep(5000 * time.Millisecond)
	lastNode := nw
	//each node joins by doing a lookup on the first node and populating its own table
	h := 1
	for k, v := range IDRTList {
		if k != *firstNode.ID {

			fmt.Println("Ny Nod varv " + strconv.Itoa(h) + ": " + v.node.rt.me.String())
			kademlia := NewKademlia(v)
			//Add first contact node to own RT
			v.node.rt.AddContact(firstNodeRT.me)

			//Do lookup on own id
			lookupResult := kademlia.LookupContact(IDRTList[k].node.rt.me.ID, IDRTList)

			//Add results from lookup to own RT
			for q := range lookupResult {
				//fmt.Println(lookupResult[q].String())
				v.node.rt.AddContact(lookupResult[q])
			}
		}
		time.Sleep(500 * time.Millisecond)
		lastNode = v
		h++	
	}

	//print the table of all nodes
	
	/*for q, w := range IDRTList {
		fmt.Println("Node: " + q.String())
		for z := range w.buckets {
			contactList := w.buckets[z]
			fmt.Println("Bucket: " + strconv.Itoa(z))
			for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
				contact := elt.Value.(Contact)
				fmt.Println(contact.String())
			}
		}
	}*/
	
	
	//print the table of the first node
	/*fmt.Println("Node: " + firstNode.ID.String())
	for z := range firstNodeRT.buckets {
		contactList := firstNodeRT.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}

	//print the table of the first node
	fmt.Println("Node: " + lastNode.me.ID.String())
	for z := range lastNode.buckets {
		contactList := lastNode.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}*/

	kademlia := NewKademlia(lastNode)
	kademlia.Store(NewKademliaID("FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "data to store", IDRTList)
	for k1, v := range IDRTList {
		for k2, v2 := range v.node.data {
			fmt.Println("Node " + k1.String() + " has " + v2 + " stored for key " + k2.String())
		}
	}

	/*data := &Data{[]byte("test")}
	// ...

	// Write the new address book back to disk.
	out, err := proto.Marshal(data)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		fmt.Println(out)
	}

	newData := &Data{}

	errUnmarsh := proto.Unmarshal(out, newData)
	if err != nil {
		log.Fatal("unmarshaling error: ", errUnmarsh)
	}
	// newData now holds {data:"test"}
	fmt.Println(newData)*/
	/*
	go Listen("localhost", 8000)
	netw := NewNetwork(lastNode)
	netw.SendPingMessage(&firstNode)
	*/
	//netw.SendFindContactMessage(&firstNode)
	//netw.SendFindDataMessage(&firstNode)

}