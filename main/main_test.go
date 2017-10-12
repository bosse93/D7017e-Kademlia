package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"os"
	"strconv"
	"time"
)

var network *Network = CreateTestNodes(20)

func connect(Usage string, arg0 string) {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Fprintf(conn, Usage+arg0)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}

	conn.Close()
}

func TestNewKademliaID(t *testing.T) {
	f := "FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	NewKademliaID(f)
}

func TestHashKademliaID(t *testing.T) {
	HashKademliaID("testar")
}

func CreateTestNodes(amount int) (network *Network)  {
	firstNode := NewContact(HashKademliaID("0"), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	node := NewNode(firstNodeRT)
	lastTCPNetwork := NewFileNetwork(node, "localhost", 8000)
	network = NewNetwork(node, lastTCPNetwork, "localhost", 8000)
	//nodeList := []*RoutingTable{firstNodeRT}
	//lastNode := firstNode
	//create 100 nodes
	if _, err := os.Stat("kademliastorage/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("kademliastorage/"+firstNode.ID.String(), 0777)
	}

	if _, err := os.Stat("upload/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("upload/"+firstNode.ID.String(), 0777)
	}

	if _, err := os.Stat("downloads/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("downloads/"+firstNode.ID.String(), 0777)
	}

	for i := 1; i < amount; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)

		ID := HashKademliaID(strconv.Itoa(i))
		rt := NewRoutingTable(NewContact(ID, a))
		//nodeList = append(nodeList, rt)
		rt.AddContact(firstNodeRT.me)
		node := NewNode(rt)
		tcpNetwork := NewFileNetwork(node, "localhost", port)
		nw := NewNetwork(node, tcpNetwork, "localhost", port)
		fmt.Println("Ny Nod varv " + strconv.Itoa(i+1) + ": " + rt.me.String())
		//go nw.Listen("localhost", port)
		time.Sleep(500 * time.Millisecond)
		kademlia := NewKademlia(nw)

		contactResult, _ := kademlia.LookupContact(ID, false)
		if len(contactResult) > 0 {
			for q := range contactResult {
				rt.AddContact(contactResult[q])
			}
		}

		if _, err := os.Stat("kademliastorage/" + ID.String()); os.IsNotExist(err) {
			os.Mkdir("kademliastorage/"+ID.String(), 0777)
		}

		if _, err := os.Stat("upload/" + ID.String()); os.IsNotExist(err) {
			os.Mkdir("upload/"+ID.String(), 0777)
		}

		if _, err := os.Stat("downloads/" + ID.String()); os.IsNotExist(err) {
			os.Mkdir("downloads/"+ID.String(), 0777)
		}

	}
	return
}

/* TODO - GÖR FRONTEND ANROPET AUTOMAGISKT från funktionen, dvs samma sak som dfs store gör */

//oklart hur detta blir, körde en dfs store från front och fick 100%
/*func TestHandleRequest(t *testing.T) {

	p := make([]byte, 2048)
	var mutex = &sync.Mutex{}
	pinned := make(map[string]bool)
	addr := net.UDPAddr{
		Port: 1234,
		IP: net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	_,remoteaddr,err := ser.ReadFromUDP(p)
	fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
	n := bytes.IndexByte(p, 0)
	split := strings.Split(string(p[:n]), " ")
	if err !=  nil {
		fmt.Printf("Some error  %v", err)
	}
	//go sendResponse(ser, remoteaddr)
	go HandleRequest(ser, remoteaddr, split, network, &pinned, mutex)
	time.Sleep(3*time.Second)
	//connect("Store", "hej")
}

func TestStartNetwork(t *testing.T) {

}


//Needs network of nodes to test this - cant do it in kademlia test without rewriting createNodes
func TestNewKademlia(t *testing.T) {
	x := NewKademlia(network)
	x = x
}*/
