package main

import (
	"net"
	"log"
	"D7024e-Kademlia/protobuf/proto"
	"fmt"
	"strconv"
	"time"
	//"math/rand"
	"sync"
)

type Network struct {
	node *Node
	waitingAnswerList map[KademliaID](chan *WrapperMessage)
	listenConnection *net.UDPConn
	threadChannels [](chan string)
	mux *sync.Mutex
}

type StringContact struct {
	ID string
	Address string
	distance string
}


func NewStringContact(id string, address string, distance string) StringContact{
	return StringContact{id, address, distance}
}

func NewNetwork(node *Node) *Network {
	network := &Network{}
	network.node = node
	network.waitingAnswerList = make(map[KademliaID]chan *WrapperMessage)
	network.mux = &sync.Mutex{}
	return network
}

//Listening for new packets on ip, port combination
func (network *Network) Listen(ip string, port int) {
	// ESTABLISH UDP CONNECTION
	serverAddr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(port))
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	network.listenConnection = serverConn

	defer serverConn.Close()

	buf := make([]byte, 4096)
	fmt.Println("Listening on port " + strconv.Itoa(port))

	//For each new packet do marshalling
	for {
		n, addr, _ := serverConn.ReadFromUDP(buf)
		wrapperRequest := &WrapperMessage{}
		replyErr := proto.Unmarshal(buf[0:n], wrapperRequest)

		go network.handleRequest(wrapperRequest, replyErr, addr)
		/*

		if wrapperRequest.Id == "ping" && replyErr == nil {
			
		
			
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM1().Id + " from " + addr.String())
		

			packet := &Reply{wrapperRequest.GetM1().Id, serverAddr.String()}
			wrapperMsg := &WrapperMessage_M4{packet}
			wrapper := &WrapperMessage{"reply", wrapperMsg}

			data, err := proto.Marshal(wrapper)
			if err != nil {
				log.Fatal("marshalling error: ", err)
			}

			buf := []byte(data)
			_, err = serverConn.WriteToUDP(buf, addr)
			if err != nil {
				log.Println(err)
			}

			
			//conn.Close()
			

		} else if wrapperRequest.Id == "contact" && replyErr == nil {
			
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM2().Id + " from " + addr.String())

		} else if wrapperRequest.Id == "data" && replyErr == nil {
			
			fmt.Println("Recieved request packet with " + wrapperRequest.Id + ", id:" + wrapperRequest.GetM3().Id + " from " + addr.String())

		} else if wrapperRequest.Id == "store" && replyErr == nil {
			

		} else if wrapperRequest.Id == "reply" && replyErr == nil {
			
			requestID, err := strconv.Atoi(wrapperRequest.GetM4().GetId())

			if err != nil {
				fmt.Println("Error")
			}
			answerChannel := network.waitingAnswerList[requestID]

			answerChannel <- wrapperRequest.GetM4().GetData()
			//close(answerChannel)

		} else {
			log.Println("Something went wrong in Listen, err: ", replyErr)
		}


		if err != nil {
			log.Fatal("Error: ", err)
		}*/
	}


}


func (network *Network) SendPingMessage(contact *Contact) Contact{
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestPing{messageID.String()}
	wrapperMsg := &WrapperMessage_M1{packet}
	wrapper := &WrapperMessage{"ping", network.rt.me.ID.String(),wrapperMsg}
	wrapper := &WrapperMessage{"ping", network.node.rt.me.ID.String(), wrapperMsg}
	
	answerChannel := make(chan *WrapperMessage)
	network.AddToChannelMap(*messageID, answerChannel)

	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)
	wrapper = network.waitForAnswer(answerChannel)	

	if(wrapper != nil) {
		returnContact := NewContact(NewKademliaID(wrapper.GetM5().Contacts[0].ID), wrapper.GetM5().Contacts[0].Address)
		return returnContact
	} else {
		fmt.Println("Timeout")
		return NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "0.0.0.0:0000")
	}
}

func (network *Network) waitForAnswer(answerChannel chan *WrapperMessage) *WrapperMessage{
	timeout := make(chan bool, 1)
	go TimeoutWaiter(5, timeout)
	for{
		select{
			case answer := <-answerChannel:
				return answer
			case <- timeout:
				return nil
		}
	}
}

func (network *Network) AddToChannelMap(messageID KademliaID, answerChannel chan *WrapperMessage) {
	network.mux.Lock()
	network.waitingAnswerList[messageID] = answerChannel

	network.mux.Unlock()
}

func TimeoutWaiter(sleepTime int, sendChannel chan bool) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
	sendChannel <-true
	close(sendChannel)
}



func (network *Network) SendFindContactMessage(contact *Contact, targetID *KademliaID) []Contact{
	// ESTABLISH UDP CONNECTION
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestContact{messageID.String(), targetID.String()}  //EDIT ME
	wrapperMsg := &WrapperMessage_M2{packet}
	wrapper := &WrapperMessage{"RequestContact", network.rt.me.ID.String(), wrapperMsg}
	wrapper := &WrapperMessage{"RequestContact", network.node.rt.me.ID.String(), wrapperMsg}

	answerChannel := make(chan *WrapperMessage, 1)
	network.AddToChannelMap(*messageID, answerChannel)
	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)

	wrapper = network.waitForAnswer(answerChannel)

	contactList := []Contact{}
	if(wrapper != nil) {
		for i := range wrapper.GetM5().GetContacts() {
			contactList = append(contactList, NewContact(NewKademliaID(wrapper.GetM5().Contacts[i].GetID()), wrapper.GetM5().Contacts[i].GetAddress()))
		}
		return contactList
	} else {
		fmt.Println("Timeout")
		contactList = append(contactList, NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "0.0.0.0:0000"))
		return contactList
	}

}

func (network *Network) SendFindDataMessage(hash string) {

}

func (network *Network) SendStoreMessage(hash string, data string, address string) {
	fmt.Println("Sending store message")
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	CheckError(err)
	sendData := []*ReplyContact_Contact{}
	sendData = append(sendData, &ReplyContact_Contact{hash, data, ""})
	packet := &ReplyContact{messageID.String(), sendData}  //EDIT ME
	wrapperMsg := &WrapperMessage_M5{packet}
	wrapper := &WrapperMessage{"store", network.node.rt.me.ID.String(), wrapperMsg}

	answerChannel := make(chan *WrapperMessage, 1)
	network.AddToChannelMap(*messageID, answerChannel)

	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)

	wrapper = network.waitForAnswer(answerChannel)
}





func (network *Network) marshalHelper(wrapper *WrapperMessage) []byte{
	data, err := proto.Marshal(wrapper)
	if err != nil {
		log.Fatal("Marshall Error: ", err)
	}
	return data
}

func (network *Network) sendPacket(data []byte, targetAddress *net.UDPAddr) {
	buf := []byte(data)
	_, err := network.listenConnection.WriteToUDP(buf, targetAddress)
	if err != nil {
		log.Println(err)
	}
}

func (network *Network) handleRequest(message *WrapperMessage, replyErr error, sourceAddress *net.UDPAddr) {
	if message.Id == "RequestPing" && replyErr == nil {	

		contakter := &ReplyContact_Contact{network.node.rt.me.ID.String(), network.node.rt.me.Address, network.node.rt.me.distance.String()}
		kontakter := []*ReplyContact_Contact{contakter}
		packet := &ReplyContact{message.GetM1().Id, kontakter}

		wrapperMsg := &WrapperMessage_M5{packet}
		wrapper := &WrapperMessage{"ReplyContact", network.rt.me.ID.String(), wrapperMsg}
		wrapper := &WrapperMessage{"ReplyContact", network.node.rt.me.ID.String(), wrapperMsg}

		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)
		
		
	} else if message.Id == "RequestContact" && replyErr == nil {
		closestContacts := network.rt.FindClosestContacts(NewKademliaID(message.GetM2().Target), 20)
		network.rt.AddContact(NewContact(NewKademliaID(message.SourceID), sourceAddress.String()))

		closestContacts := network.node.rt.FindClosestContacts(NewKademliaID(message.GetM2().Target), 20)
		kontakter := []*ReplyContact_Contact{}
		for i := range closestContacts {
			contakter := &ReplyContact_Contact{closestContacts[i].ID.String(), closestContacts[i].Address, closestContacts[i].String()}
			kontakter = append(kontakter, contakter)
		}

		packet := &ReplyContact{message.GetM2().GetId(), kontakter}
		wrapperMsg := &WrapperMessage_M5{packet}
		wrapper := &WrapperMessage{"ReplyContact", network.rt.me.ID.String(), wrapperMsg}
		wrapper := &WrapperMessage{"ReplyContact", network.node.rt.me.ID.String(), wrapperMsg}

		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)

	} else if message.Id == "data" && replyErr == nil {
		

	} else if message.Id == "store" && replyErr == nil {
		fmt.Println("Received store")
		//store data (string) in data map
		network.node.Store(NewKademliaID(message.GetM5().Contacts[0].ID), message.GetM5().Contacts[0].Address)
		//send reply
		packet := &Reply{message.SourceID, "ok"}
		wrapperMsg := &WrapperMessage_M4{packet}
		wrapper := &WrapperMessage{"Reply", network.node.rt.me.ID.String(), wrapperMsg}
		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)

	} else if message.Id == "Reply" && replyErr == nil {
		fmt.Println("Got reply " + message.GetM4().Data)
	} else if message.Id == "ReplyContact" && replyErr == nil {
		network.mux.Lock()
		requestID := NewKademliaID(message.GetM5().GetId())

		answerChannel := network.waitingAnswerList[*requestID]

		if(answerChannel != nil) {
			answerChannel <- message
		} else {
			fmt.Println("Forged Reply")
		}
		
		close(answerChannel)
		network.mux.Unlock()

	} else {
		fmt.Println(message.Id)
		log.Println("Something went wrong in Listen, err: ", replyErr)
	}

}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}