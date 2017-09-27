package main

import (
	"fmt"
	"strconv"
	"time"
	"net"
)

func main() {
	startNetwork()
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func handleResponse(conn *net.UDPConn, addr *net.UDPAddr, p string){
	_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage " + p), addr)

	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func startNetwork() {
	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	lastNetwork := NewNetwork(NewNode(firstNodeRT), "localhost", 8000)

	nodeList := []*RoutingTable{firstNodeRT}
	//lastNode := firstNode
	//create 100 nodes
	for i := 0; i < 10; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)


		ID := NewRandomKademliaID()
		rt := NewRoutingTable(NewContact(ID, a))
		nodeList = append(nodeList, rt)
		rt.AddContact(firstNodeRT.me)
		nw := NewNetwork(NewNode(rt), "localhost", port)
		fmt.Println("Ny Nod varv " + strconv.Itoa(i+1) + ": " + rt.me.String())
		//go nw.Listen("localhost", port)
		time.Sleep(500 * time.Millisecond)
		kademlia := NewKademlia(nw)

		contactResult, _  := kademlia.LookupContact(ID, false)
		if(len(contactResult) > 0) {
			for q := range contactResult {
				rt.AddContact(contactResult[q])
			}
		}
		//lastNetwork = nw
	}

	printFirstNodeRT(firstNode, firstNodeRT)
	printLastNodeRT(nodeList)


	kademlia := NewKademlia(lastNetwork)
	kademlia.Store(NewKademliaID("FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "data to store")
	time.Sleep(3*time.Second)
	kademlia = NewKademlia(lastNetwork)
	data, success := kademlia.LookupData("FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	if(success) {
		fmt.Println("Data returned " + data)
	} else {
		fmt.Println("Data not found")
	}



	//TEST
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP: net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_,remoteaddr,err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		fmt.Println(p)
		fmt.Println(string(p))
		if string(p) =="Cat"{
			fmt.Println("I got cat back do cat stuff")
		} else if string(p) == "Store" {
			fmt.Println("I got Store back")
		}
		if err !=  nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		//go sendResponse(ser, remoteaddr)
		go handleResponse(ser, remoteaddr, string(p))
	}

	/*for k1, v := range IDRTList {
		for k2, v2 := range v.node.data {
			fmt.Println("Node " + k1.String() + " has " + v2 + " stored for key " + k2.String())
		}
	}*/

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
