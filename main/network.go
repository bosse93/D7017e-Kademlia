package main

import (
	"net"
	"log"
	"D7024e-Kademlia/protobuf/proto"
	"fmt"
	"strconv"
	"time"
	"sync"
)

type Network struct {
	node *Node
	waitingAnswerList map[KademliaID](chan *WrapperMessage)
	returnDataChannels map[KademliaID](*chan string)
	listenConnection *net.UDPConn
	threadChannels [](chan string)
	mux *sync.Mutex
}

func NewNetwork(node *Node, ip string, port int) *Network {
	network := &Network{}
	network.node = node
	network.waitingAnswerList = make(map[KademliaID]chan *WrapperMessage)
	network.mux = &sync.Mutex{}

	// ESTABLISH UDP CONNECTION
	serverAddr, err := net.ResolveUDPAddr("udp", ip + ":" + strconv.Itoa(port))
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	network.listenConnection = serverConn
	buf := make([]byte, 4096)
	fmt.Println("Listening on port " + strconv.Itoa(port))
	go network.Listen(buf)

	return network
}

//Listening for new packets on ip, port combination
func (network *Network) Listen(buf []byte) {
	defer network.listenConnection.Close()

	//For each new packet do marshalling
	for {
		n, addr, _ := network.listenConnection.ReadFromUDP(buf)
		wrapperRequest := &WrapperMessage{}
		replyErr := proto.Unmarshal(buf[0:n], wrapperRequest)

		go network.handleRequest(wrapperRequest, replyErr, addr)
	}
}


func (network *Network) SendPingMessage(contact *Contact) Contact{
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestPing{messageID.String()}
	wrapperMsg := &WrapperMessage_M1{packet}
	wrapper := &WrapperMessage{"ping", network.node.rt.me.ID.String(), wrapperMsg}
	
	answerChannel := network.createChannel(*messageID)
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

func (network *Network) SendFindContactMessage(contact *Contact, targetID *KademliaID) []Contact{
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestContact{messageID.String(), targetID.String()}
	wrapperMsg := &WrapperMessage_M2{packet}
	wrapper := &WrapperMessage{"RequestContact", network.node.rt.me.ID.String(), wrapperMsg}

	answerChannel := network.createChannel(*messageID)
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


func (network *Network) SendFindDataMessage(hash string, contact Contact) []Contact {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestData{messageID.String(), hash}
	wrapperMsg := &WrapperMessage_M3{packet}
	wrapper := &WrapperMessage{"RequestData", network.node.rt.me.ID.String(), wrapperMsg}

	answerChannel := network.createChannel(*messageID)
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

	answerChannel := network.createChannel(*messageID)

	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)

	wrapper = network.waitForAnswer(answerChannel)
}

func (network *Network) handleRequest(message *WrapperMessage, replyErr error, sourceAddress *net.UDPAddr) {
	if message.Id == "RequestPing" && replyErr == nil {	

		contactReply := &ReplyContact_Contact{network.node.rt.me.ID.String(), network.node.rt.me.Address, network.node.rt.me.distance.String()}
		contactListReply := []*ReplyContact_Contact{contactReply}
		packet := &ReplyContact{message.GetM1().Id, contactListReply}
		wrapperMsg := &WrapperMessage_M5{packet}
		wrapper := &WrapperMessage{"ReplyContact", network.node.rt.me.ID.String(), wrapperMsg}

		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)
		
		
	} else if message.Id == "RequestContact" && replyErr == nil {
		closestContacts := network.node.rt.FindClosestContacts(NewKademliaID(message.GetM2().Target), 20)
		network.node.rt.AddContact(NewContact(NewKademliaID(message.SourceID), sourceAddress.String()))

		contactListReply := []*ReplyContact_Contact{}
		for i := range closestContacts {
			contactReply := &ReplyContact_Contact{closestContacts[i].ID.String(), closestContacts[i].Address, closestContacts[i].String()}
			contactListReply = append(contactListReply, contactReply)
		}
		packet := &ReplyContact{message.GetM2().GetId(), contactListReply}
		wrapperMsg := &WrapperMessage_M5{packet}
		wrapper := &WrapperMessage{"ReplyContact", network.node.rt.me.ID.String(), wrapperMsg}

		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)

	} else if message.Id == "RequestData" && replyErr == nil {
		packet := &ReplyData{}
		network.node.rt.AddContact(NewContact(NewKademliaID(message.SourceID), sourceAddress.String()))
		packet.Id = message.GetM3().Id
		if val, ok := network.node.data[NewKademliaID(message.GetM3().Key)]; ok {
			packet.ReturnType = "data"
			reply := &Reply{message.GetM3().Key, val}
			dataPacket := &ReplyData_ReplyData{reply}
			packet.Msg = dataPacket
		} else {
			packet.ReturnType = "contact"
			closestContacts := network.node.rt.FindClosestContacts(NewKademliaID(message.GetM3().Key), 20)
			contactListReply := []*ReplyContact_Contact{}
			for i := range closestContacts {
				contactReply := &ReplyContact_Contact{closestContacts[i].ID.String(), closestContacts[i].Address, closestContacts[i].String()}
				contactListReply = append(contactListReply, contactReply)
			}
			reply := &ReplyContact{message.GetM3().Key, contactListReply}
			contactPacket := &ReplyData_ReplyContact{reply}
			packet.Msg = contactPacket
		}
		wrapperMsg := &WrapperMessage_ReplyData{packet}
		wrapper := &WrapperMessage{"ReplyData", network.node.rt.me.ID.String(), wrapperMsg}

		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)

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
		requestID := NewKademliaID(message.GetM5().GetId())

		network.mux.Lock()
		answerChannel := network.waitingAnswerList[*requestID]
		network.mux.Unlock()

		if(answerChannel != nil) {
			answerChannel <- message
		} else {
			fmt.Println("Forged Reply")
		}
		close(answerChannel)

	} else if message.Id == "ReplyData" {
		requestID := NewKademliaID(message.GetReplyData().GetId())

		network.mux.Lock()
		answerChannel := network.waitingAnswerList[*requestID]
		network.mux.Unlock()

		if message.GetReplyData().ReturnType == "data" {
			returnChannel := *network.returnDataChannels[*NewKademliaID(message.GetReplyData().GetReplyData().Id)]
			returnChannel <- message.GetReplyData().GetReplyData().Data
		}
		if answerChannel != nil {
			answerChannel <- message
		} else {
			fmt.Println("Forged Reply")
		}
		close(answerChannel)
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

func (network *Network) createChannel(messageID KademliaID) chan *WrapperMessage{
	answerChannel := make(chan *WrapperMessage, 1)
	network.mux.Lock()
	network.waitingAnswerList[messageID] = answerChannel
	network.mux.Unlock()
	return answerChannel
}

func TimeoutWaiter(sleepTime int, sendChannel chan bool) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
	sendChannel <-true
	close(sendChannel)
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

func (network *Network) setupDataChannel(id KademliaID, c *chan string) {
	network.returnDataChannels[id] = c
}

