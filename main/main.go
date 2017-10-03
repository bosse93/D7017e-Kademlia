package main

import (
	"fmt"
	"strconv"
	"time"
	"net"
	"encoding/hex"
	"io/ioutil"
	//"os"
)

func main() {
	StartNetwork()
}

func HashKademliaID(fileName string) *KademliaID{
	f := hex.EncodeToString([]byte(fileName))
	for len(f) < 40 {
		f = f + "0"
	}
	return NewKademliaID(f)
}

func HandleRequest(conn *net.UDPConn, addr *net.UDPAddr, p string, network *Network){
	//_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage " + p), addr)

	if p[:5]=="Store" {
		fmt.Println("this was a store message with arg "+ p[5:])
		kademlia := NewKademlia(network)
		//FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

		newKad := HashKademliaID(p[5:])
		kademlia.Store(newKad, "data to store")
		_,storeErr := conn.WriteToUDP([]byte("stored: "+newKad.String()), addr)
		if storeErr != nil {
			fmt.Println("something went shit in store: %v", storeErr)
		}

	} else if p[:3]=="Cat" {
		fmt.Println("I got a Cat call")
	}

}

func CreateNodes(amount int) *Network{
	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	lastNetwork := NewNetwork(NewNode(firstNodeRT), "localhost", 8000)

	nodeList := []*RoutingTable{firstNodeRT}
	//lastNode := firstNode
	//create 100 nodes
	for i := 0; i < amount; i++ {
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
		lastNetwork = nw

	}
	return lastNetwork
}

func StartFrontend(lastNetwork *Network){
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
		go HandleRequest(ser, remoteaddr, string(p), lastNetwork)

	}
}

func StartNetwork() {
	//Creates x amount of nodes in a network
	lastNetwork := CreateNodes(10)

	//printFirstNodeRT(firstNode, firstNodeRT)
	//printLastNodeRT(nodeList)

	testStore := HashKademliaID("testStore.txt")
	fmt.Println("test store " + testStore.String())

	//pwd, _ := os.Getwd()
	//dat, err := ioutil.ReadFile(pwd+"/../src/D7024e-Kademlia/main/testStore.txt")
	dat, err := ioutil.ReadFile("main/testStore.txt")
	check(err)

	kademlia := NewKademlia(lastNetwork)
	kademlia.Store(testStore, string(dat))
	time.Sleep(3*time.Second)
	kademlia = NewKademlia(lastNetwork)
	data, success := kademlia.LookupData(testStore.String())
	if(success) {
		fmt.Println("Data returned " + data)
	} else {
		fmt.Println("Data not found")
	}

	//Setup Frontend
	StartFrontend(lastNetwork)


	/*for k1, v := range IDRTList {
		for k2, v2 := range v.node.data {
			fmt.Println("Node " + k1.String() + " has " + v2 + " stored for key " + k2.String())
		}
	}*/

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



func check(e error) {
	if e != nil {
		panic(e)
	}
}