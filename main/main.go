package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	StartNetwork()
}

func HandleRequest(conn *net.UDPConn, addr *net.UDPAddr, args []string, network *Network, pinned *map[string]bool, mux *sync.Mutex) {

	if args[0] == "store" {
		fmt.Println("this was a store message with arg " + args[0])
		kademlia := NewKademlia(network)
		//FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

		newKad := HashKademliaID(args[1])
		upload(network.node.rt.me.ID.String(), args[1])
		kademlia.Store(args[1])
		_, storeErr := conn.WriteToUDP([]byte("stored: "+newKad.String()), addr)
		if storeErr != nil {
			fmt.Println("something went shit in store: %v", storeErr)
		}

	} else if args[0] == "cat" {
		kademlia := NewKademlia(network)
		//FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

		success := kademlia.LookupData(args[1])
		if success {
			mux.Lock()
			(*pinned)[args[1]] = false
			mux.Unlock()
			go RemoveFile(20, pinned, args[1], network.node.rt.me.ID.String(), mux)
			time.Sleep(1000 * time.Millisecond)
			dat, readerr := ioutil.ReadFile("downloads/" + network.node.rt.me.ID.String() + "/" + args[1])
			if readerr != nil {
				log.Fatal(readerr)
			}
			fmt.Println("file content: " + string(dat))
			_, err := conn.WriteToUDP([]byte(dat), addr)
			if err != nil {
				fmt.Println("something went shit in lookup: %v", err)
			}
		} else {
			_, err := conn.WriteToUDP([]byte("no data found"), addr)
			if err != nil {
				fmt.Println("something went shit in lookup: %v", err)
			}
		}
	} else if args[0] == "pin" {
		(*pinned)[args[1]] = true
	} else if args[0] == "unpin" {
		(*pinned)[args[1]] = false
	}

}

func CreateNodes(amount int) (firstNetwork *Network) {

	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	node := NewNode(firstNodeRT)
	lastTCPNetwork := NewFileNetwork(node, "localhost", 8000)
	firstNetwork = NewNetwork(node, lastTCPNetwork, "localhost", 8000)

	if _, err := os.Stat("kademliastorage/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("kademliastorage/"+firstNode.ID.String(), 0777)
	}

	if _, err := os.Stat("upload/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("upload/"+firstNode.ID.String(), 0777)
	}

	if _, err := os.Stat("downloads/" + firstNode.ID.String()); os.IsNotExist(err) {
		os.Mkdir("downloads/"+firstNode.ID.String(), 0777)
	}

	for i := 0; i < amount; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)

		ID := NewRandomKademliaID()
		rt := NewRoutingTable(NewContact(ID, a))
		rt.AddContact(firstNodeRT.me)
		node := NewNode(rt)
		tcpNetwork := NewFileNetwork(node, "localhost", port)
		nw := NewNetwork(node, tcpNetwork, "localhost", port)
		fmt.Println("Ny Nod varv " + strconv.Itoa(i+1) + ": " + rt.me.String())
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

func StartFrontend(lastNetwork *Network) {
	var mutex = &sync.Mutex{}
	pinned := make(map[string]bool)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		p := make([]byte, 2048)
		_, remoteaddr, err := ser.ReadFromUDP(p)
		n := bytes.IndexByte(p, 0)
		split := strings.Split(string(p[:n]), " ")
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		fmt.Println(split)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go HandleRequest(ser, remoteaddr, split, lastNetwork, &pinned, mutex)
		time.Sleep(100 * time.Millisecond)
		for key, value := range pinned {
			fmt.Println("key: " + key + ", value:" + strconv.FormatBool(value))
		}
	}
}

func StartNetwork() {
	if _, err := os.Stat("kademliastorage/"); os.IsNotExist(err) {
		os.Mkdir("kademliastorage", 0777)
	} else {
		os.RemoveAll("kademliastorage")
		time.Sleep(500 * time.Millisecond)
		os.Mkdir("kademliastorage", 0777)
	}

	if _, err := os.Stat("upload/"); os.IsNotExist(err) {
		os.Mkdir("upload", 0777)
	} else {
		os.RemoveAll("upload")
		time.Sleep(500 * time.Millisecond)
		os.Mkdir("upload", 0777)
	}

	if _, err := os.Stat("downloads/"); os.IsNotExist(err) {
		os.Mkdir("downloads", 0777)
	}
	firstNetwork := CreateNodes(100)
	StartFrontend(firstNetwork)
}

func upload(id string, file string) {
	fileDst, dstErr := os.Create("upload/" + id + "/" + file)
	if dstErr != nil {
		log.Fatal(dstErr)
	} else {
		fmt.Println(fileDst)
	}
	fileSrc, fileErr := os.Open(file)
	if fileErr != nil {
		log.Fatal(fileErr)
	} else {
		fmt.Println("fileSrc ok")
	}
	if _, err := io.Copy(fileDst, fileSrc); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("fileDst ok")
	}
	fileSrc.Close()
	fileDst.Close()
}

func RemoveFile(sleepTime int, pinned *map[string]bool, file string, id string, mux *sync.Mutex) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
	fmt.Println("removing " + "downloads/" + id + "/" + file + " if not pinned")

	if _, err := os.Stat("downloads/" + id + "/" + file); !os.IsNotExist(err) {
		if !(*pinned)[file] {
			fmt.Println("not pinned")
			mux.Lock()
			os.Remove("downloads/" + id + "/" + file)
			mux.Unlock()
		} else {
			fmt.Println("pinned, trying again later")
			go RemoveFile(sleepTime, pinned, file, id, mux)
		}
	}

}
