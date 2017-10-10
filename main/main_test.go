package main

import (
	"testing"
	"net"
	"fmt"
	"bufio"
	"time"
	"sync"
	"strings"
	"bytes"
)

var network *Network

func connect(Usage string, arg0 string){
	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Fprintf(conn, Usage + arg0)
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

func TestCreateNodes(t *testing.T) {
	//CreateNodes(10)
	network = CreateNodes(10)
}

/* TODO - GÖR FRONTEND ANROPET AUTOMAGISKT från funktionen, dvs samma sak som dfs store gör */
//oklart hur detta blir, körde en dfs store från front och fick 100%
func TestStartFrontend(t *testing.T) {
	go StartFrontend(network)
	time.Sleep(3*time.Second)
	connect("Store", "hej")
}


//oklart hur detta blir, körde en dfs store från front och fick 100%
func TestHandleRequest(t *testing.T) {

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
}