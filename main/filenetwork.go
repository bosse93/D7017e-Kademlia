package main

import (
	"net"
	"log"
	"fmt"
	"strconv"
	"os"
	"strings"
	"io"
)

type FileNetwork struct {
	node *Node
	//waitingAnswerList map[KademliaID](chan interface{})
	listenConnection *net.TCPListener
}


func NewFileNetwork(node *Node, ip string, port int) *FileNetwork {
	network := &FileNetwork{}
	network.node = node

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
	 buffer := make([]byte, 1024)

	 _, error := connection.Read(buffer)
    if error != nil {
        fmt.Println("There is an error reading from connection", error.Error())
        return
    }

    file, err := os.Open(strings.TrimSpace(string(buffer))) // For read access.
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close() // make sure to close the file even if we panic.
	n, err := io.Copy(connection, file)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(n, "bytes sent")


}

func (network *FileNetwork) downloadFile(fileID *KademliaID, address string) {
	destinationAddr, err := net.ResolveTCPAddr("tcp", address)
	connection, err := net.DialTCP("tcp", nil, destinationAddr)
    if err != nil {
        fmt.Println("There was an error making a connection")
    }

	//var currentByte int64 = 0

    //fileBuffer := make([]byte, 1024)

    //var err error
    file, err := os.Create(strings.TrimSpace(fileID.String()))
    if err != nil {
        log.Fatal(err)
    }
    connection.Write([]byte(fileID.String()))

    defer file.Close() // make sure to close the file even if we panic.
	n, err := io.Copy(file, connection)
	if err != nil {
    	log.Fatal(err)
	}
	fmt.Println(n, "bytes copied")
 

    return


}