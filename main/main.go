package main

import (
	"fmt"
	"strconv"
	"time"
	"net"
	"encoding/hex"
	//"io/ioutil"
	"strings"
	"os"
	"net/http"
	"io"
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

func HandleRequest(conn *net.UDPConn, addr *net.UDPAddr, args []string, network *Network){
	//_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage " + p), addr)

	if args[0]=="store" {
		fmt.Println("this was a store message with arg "+ args[0])
		kademlia := NewKademlia(network)
		//FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

		newKad := HashKademliaID(args[1])
		kademlia.Store(newKad, args[2])
		_,storeErr := conn.WriteToUDP([]byte("stored: "+newKad.String()), addr)
		if storeErr != nil {
			fmt.Println("something went shit in store: %v", storeErr)
		}

	} else if args[0]=="cat" {
		kademlia := NewKademlia(network)
		//FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

		newKad := HashKademliaID(args[1])
		data, success := kademlia.LookupData(newKad.String())
		if success {
			_,err := conn.WriteToUDP([]byte(data), addr)
			if err != nil {
				fmt.Println("something went shit in lookup: %v", err)
			}
		} else {
			_,err := conn.WriteToUDP([]byte("no data found"), addr)
			if err != nil {
				fmt.Println("something went shit in lookup: %v", err)
			}
		}

	}

}

func CreateNodes(amount int) *Network{
	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	node := NewNode(firstNodeRT)
	lastTCPNetwork := NewFileNetwork(node, "localhost", 8000)
	lastNetwork := NewNetwork(node, lastTCPNetwork, "localhost", 8000)


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
		node := NewNode(rt)
		tcpNetwork := NewFileNetwork(node, "localhost", port)
		nw := NewNetwork(node, tcpNetwork, "localhost", port)
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
		p := make([]byte, 2048)
		_,remoteaddr,err := ser.ReadFromUDP(p)
		split := strings.Split(string(p), " ")
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		fmt.Println(split)
		if split[0] == "cat"{
			fmt.Println("I got cat back do cat stuff")
		} else if split[0] == "store" {
			fmt.Println("I got Store back")
		}
		if err !=  nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		//go sendResponse(ser, remoteaddr)
		go HandleRequest(ser, remoteaddr, split, lastNetwork)

	}
}

func StartNetwork() {
	//Creates x amount of nodes in a network
	lastNetwork := CreateNodes(50)

	//printFirstNodeRT(firstNode, firstNodeRT)
	//printLastNodeRT(nodeList)

	testStore := HashKademliaID("workshop.jpeg")

	//pwd, _ := os.Getwd()
	//dat, err := ioutil.ReadFile(pwd+"/../src/D7024e-Kademlia/main/testStore.txt")
	//dat, err := ioutil.ReadFile("main/testStore.txt")
	//check(err)

	kademlia := NewKademlia(lastNetwork)
	//store link to workshop jpg

	kademlia.Store(testStore, "https://www.dropbox.com/s/b0a98iiuu1o9m5y/Workshopmockup-1.jpg?dl=1")
	time.Sleep(3*time.Second)
	//kademlia = NewKademlia(lastNetwork)
	//lookup workshop jpg
	/*
	data, success := kademlia.LookupData(testStore.String())
	if(success) {
		fmt.Println("Data returned " + data)
	} else {
		fmt.Println("Data not found")
	}

	//download workshop jpg, to be done in frontend when response with url arrives.
	downerr := downloadFile("workshop.jpeg", data)
	check(downerr)*/
	//Setup Frontend

	//downloadFile("workshop.jpeg", "https://www.dropbox.com/s/b0a98iiuu1o9m5y/Workshopmockup-1.jpg?dl=1")
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

func downloadFile(filepath string, url string) (err error) {
	fmt.Println("filepath: " + filepath + " url: " + url)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}