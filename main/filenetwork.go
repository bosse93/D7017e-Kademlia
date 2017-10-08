package main

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"os"
	//"strings"
	"io"
	"time"
	"sync"
)

type FileNetwork struct {
	node *Node
	listenConnection *net.TCPListener
	mux1 *sync.Mutex
	mux2 *sync.Mutex
}


func NewFileNetwork(node *Node, ip string, port int) *FileNetwork {
	network := &FileNetwork{}
	network.node = node
	network.mux1 = &sync.Mutex{}
	network.mux2 = &sync.Mutex{}

	serverAddr, err := net.ResolveTCPAddr("tcp", ip + ":" + strconv.Itoa(port))


	serverConn, err := net.ListenTCP("tcp", serverAddr)
	CheckError(err)
	network.listenConnection = serverConn
	fmt.Println("TCP Listening on port " + strconv.Itoa(port))
	go network.Listen()

	return network
}

//Listening for new packets on ip, port combination
func (network *FileNetwork) Listen() {
	defer network.listenConnection.Close()

	for {
		conn, err := network.listenConnection.Accept()
 		if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        go network.HandleFileRequest(conn)
	}
}

func (network *FileNetwork) HandleFileRequest(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 1024)

	n, error := connection.Read(buffer)
    if error != nil {
        fmt.Println("There is an error reading from connection", error.Error())
        return
    }
    //filetoOpen := string(buffer) + ".txt"

    if(network.node.gotData(*NewKademliaID(string(buffer[:n])))) {
    	network.mux2.Lock()
    	file, err := os.Open("kademliastorage/" + network.node.rt.me.ID.String() + "/" + string(buffer[:n]))
    	if err != nil {
			log.Fatal(err)
		}
		defer file.Close() // make sure to close the file even if we panic.
		n1, err := io.Copy(connection, file)
		if err != nil {
		    log.Fatal(err)
		}
		network.mux2.Unlock()
		fmt.Println(n1, "bytes sent2")
    } else {
    	if _, err := os.Stat("upload/" + network.node.rt.me.ID.String() + "/" + decodeHash(string(buffer[:n]))); os.IsNotExist(err) {
    		fmt.Println("Node doesn't have the requested file.")
    		connection.Close()
    		return
    	} else {
    		network.mux1.Lock()
    		file, err := os.Open("upload/" + network.node.rt.me.ID.String() + "/" + decodeHash(string(buffer[:n])))
    		if err != nil {
			log.Fatal(err)
		}
    		defer file.Close() // make sure to close the file even if we panic.
			n1, err := io.Copy(connection, file)
			if err != nil {
			    log.Fatal(err)
			}
			network.mux1.Unlock()
			fmt.Println(n1, "bytes sent1")
    	}
    }
}


func (network *FileNetwork) downloadFile(fileID *KademliaID, address string, userDownload bool) {
	destinationAddr, err := net.ResolveTCPAddr("tcp", address)
	connection, err := net.DialTCP("tcp", nil, destinationAddr)
	defer connection.Close()
    if err != nil {
        fmt.Println("There was an error making a connection")
    }

	//var currentByte int64 = 0
	connection.Write([]byte(fileID.String()))

    //fileBuffer := make([]byte, 1024)
    if(userDownload) {
	    file, err := os.Create("downloads/" + network.node.rt.me.ID.String() + "/" + decodeHash(fileID.String()))
	    if err != nil {
	        log.Fatal(err)
	    }

    	defer file.Close() // make sure to close the file even if we panic.
		n, err := io.Copy(file, connection)
		if err != nil {
    		log.Fatal(err)
		}
		fmt.Println(n, "bytes copied")
    } else {
    	file, err := os.Create("kademliastorage/" + network.node.rt.me.ID.String() + "/" + fileID.String())
    	if err != nil {
        	log.Fatal(err)
    	}

    	defer file.Close() // make sure to close the file even if we panic.
		n, err := io.Copy(file, connection)
		if err != nil {
    		log.Fatal(err)
		}
		network.node.Store(*fileID, time.Now())
		fmt.Println(n, "bytes copied")
    }
}