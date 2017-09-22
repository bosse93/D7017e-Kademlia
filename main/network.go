package main

import (
	"net"
	"log"
	"D7024e-Kademlia/protobuf/proto"
	"fmt"
	"strconv"
	"time"
)

type Network struct {
	rt RoutingTable
}

func NewNetwork(rt *RoutingTable) *Network {
	network := &Network{}
	network.rt = *rt
	return network
}

func Listen(ip string, port int) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	fmt.Println("server address " + serverAddr.String())
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	fmt.Println("Listening on port " + strconv.Itoa(port))
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		packetRequest := &RequestPing{}
		wrapperReply := &WrapperMessage{}


		replyErr := proto.Unmarshal(buf[0:n], wrapperReply)
		requestErr := proto.Unmarshal(buf[0:n], packetRequest)

		if wrapperReply.Id == "ping" {
			fmt.Println("Recieved reply packet with "+ wrapperReply.Id)
			//test := proto.Unmarshal(buf[0:n], wrapperReply.Msg)


		}

		if requestErr == nil {
			fmt.Println("Received request packet with " + packetRequest.Id + " from " + addr.String())
		}
		if replyErr == nil {
			//fmt.Println("Recieved reply packet with " + packetReply.Id + " and " + packetReply.Data + " from " + addr.String())
		}


		if err != nil {
			log.Fatal("Error: ", err)
		}
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	fmt.Println("remote address " + remoteAddr.String())
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", network.rt.me.Address)
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err)

	defer conn.Close()
	i := 1
	for {
		packet := &RequestPing{strconv.Itoa(i)}
		wrapperMsg := &WrapperMessage_M1{packet}
		wrapper := &WrapperMessage{"ping", wrapperMsg}

		//packet := &Reply{strconv.Itoa(i), "reply data"}
		data, err := proto.Marshal(wrapper)
		if err != nil {
			log.Fatal("marshalling error: ", err)
		}
		buf := []byte(data)
		_, err = conn.Write(buf)
		if err != nil {
			log.Println(err)
		}
		i++
		time.Sleep(time.Second * 1)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}






func (network *Network) handleRequest() {
	//Unpack packet
	//packet := ....

	/*network.rt.AddContact(packet.sender)

	switch packet.messageType {
		case "SendPingMessage":
			//Måste svara nåt! Svara med sig själv?

		case "SendFindContactMessage":
			answerData := network.rt.FindClosestContacts(packet.target, 20)

			//Pack answer. answerData, packet.messageID, ....

		case "SendFindDataMessage":
			//Check if data in storage, packet.target(ID to data)

			//If Data 
			//send file How?

			//Else
			answerData := network.rt.FindClosestContacts(packet.target, 20)

		case "SendStoreMessage":
			//Store data
			//Kolla om data redan finns?

			//Svara med något?


		default:
			//Not a valid message


	}

	//Send marshalled response on UDP*/



}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}