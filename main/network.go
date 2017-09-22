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

//Listening for new packets on ip, port combination
func Listen(ip string, port int) {
	// ESTABLISH UDP CONNECTION
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	fmt.Println("server address " + serverAddr.String())
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)
	fmt.Println("Listening on port " + strconv.Itoa(port))

	//For each new packet do marshalling
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		wrapperRequest := &WrapperMessage{}
		replyErr := proto.Unmarshal(buf[0:n], wrapperRequest)

		if wrapperRequest.Id == "ping" && replyErr == nil {
			/** We a ping **/
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM1().Id + " from " + addr.String())
			//Pinga tillbaka

		} else if wrapperRequest.Id == "contact" && replyErr == nil {
			/** We got a contact **/
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM2().Id + " from " + addr.String())

		} else if wrapperRequest.Id == "data" && replyErr == nil {
			/** We got some data **/
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM3().Id + " from " + addr.String())

		} else if wrapperRequest.Id == "store" && replyErr == nil {
			/** Store **/

		} else {
			log.Println("Something went wrong in Listen, err: ", replyErr)
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
	// ESTABLISH UDP CONNECTION
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
		packet := &RequestContact{strconv.Itoa(i)}  //EDIT ME
		wrapperMsg := &WrapperMessage_M2{packet}
		wrapper := &WrapperMessage{"contact", wrapperMsg}

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

func (network *Network) SendFindDataMessage(hash string) {

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