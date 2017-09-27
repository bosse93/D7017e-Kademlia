package main

import (
	"net"
	"log"
	"D7024e-Kademlia/github.com/protobuf/proto"
	"fmt"
	"strconv"
	"time"
	"sync"
)

type Network struct {
	node *Node
	waitingAnswerList map[KademliaID](chan interface{})
	//returnDataChannels map[KademliaID](chan string)
	listenConnection *net.UDPConn
	threadChannels [](chan string)
	mux *sync.Mutex
}

type DataReturn struct {
	contacts []Contact
	data string
}

func NewNetwork(node *Node, ip string, port int) *Network {
	network := &Network{}
	network.node = node
	//network.returnDataChannels = make(map[KademliaID]chan string)
	network.waitingAnswerList = make(map[KademliaID]chan interface{})
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

/*
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
*/

func (network *Network) SendFindContactMessage(targetID *KademliaID, contact *Contact, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestContact{messageID.String(), targetID.String()}
	wrapperMsg := &WrapperMessage_M2{packet}
	wrapper := &WrapperMessage{"RequestContact", network.node.rt.me.ID.String(), wrapperMsg}

	network.createChannel(messageID, returnChannel)
	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)
	
	go network.TimeoutWaiter(5, returnChannel, messageID)
}


func (network *Network) SendFindDataMessage(hash string, contact *Contact, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestData{messageID.String(), hash}
	wrapperMsg := &WrapperMessage_M3{packet}
	wrapper := &WrapperMessage{"RequestData", network.node.rt.me.ID.String(), wrapperMsg}

	network.createChannel(messageID, returnChannel)
	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)
	

}

func (network *Network) SendStoreMessage(hash string, data string, address string, returnChannel chan interface{}) {
	fmt.Println("Sending store message")
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	CheckError(err)
	sendData := []*ReplyContact_Contact{}
	sendData = append(sendData, &ReplyContact_Contact{hash, data, ""})
	packet := &ReplyContact{messageID.String(), sendData}  //EDIT ME
	wrapperMsg := &WrapperMessage_M5{packet}
	wrapper := &WrapperMessage{"store", network.node.rt.me.ID.String(), wrapperMsg}

	network.createChannel(messageID, returnChannel)

	network.sendPacket(network.marshalHelper(wrapper), remoteAddr)
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
		if _, ok := network.node.data[*NewKademliaID(message.GetM3().Key)]; ok {
			/*
			fmt.Println("data found")
			packet.ReturnType = "data"
			reply := &Reply{message.GetM3().Key, val}
			dataPacket := &ReplyData_ReplyData{reply}
			packet.Msg = dataPacket
			*/
		} else {
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
		}
	} else if message.Id == "store" && replyErr == nil {
		fmt.Println("Received store")
		//store data (string) in data map
		network.node.Store(*NewKademliaID(message.GetM5().Contacts[0].ID), message.GetM5().Contacts[0].Address)
		//send reply
		packet := &Reply{message.GetM5().GetId(), "ok"}
		wrapperMsg := &WrapperMessage_M4{packet}
		wrapper := &WrapperMessage{"Reply", network.node.rt.me.ID.String(), wrapperMsg}
		network.sendPacket(network.marshalHelper(wrapper), sourceAddress)

	} else if message.Id == "Reply" && replyErr == nil {
		fmt.Println("Got reply " + message.GetM4().Data)
		requestID := NewKademliaID(message.GetM4().GetId())

		network.mux.Lock()
		answerChannel := network.waitingAnswerList[*requestID]
		network.mux.Unlock()

		if(answerChannel != nil) {
			answerChannel<-message.GetM4().GetData()
		} else {
			fmt.Println("Forged Reply or Timeout")
		}

	} else if message.Id == "ReplyContact" && replyErr == nil {
		requestID := NewKademliaID(message.GetM5().GetId())

		network.mux.Lock()
		answerChannel := network.waitingAnswerList[*requestID]
		network.mux.Unlock()

		if(answerChannel != nil) {
			contactList := []Contact{}
			for i := range message.GetM5().GetContacts() {
				contactList = append(contactList, NewContact(NewKademliaID(message.GetM5().Contacts[i].GetID()), message.GetM5().Contacts[i].GetAddress()))
			}
			answerChannel<-contactList
		} else {
			fmt.Println("Forged Reply or Timeout")
		}
		
	} else if message.Id == "ReplyData" {
		fmt.Println("received data reply")
		requestID := NewKademliaID(message.GetReplyData().GetId())

		network.mux.Lock()
		answerChannel := network.waitingAnswerList[*requestID]
		network.mux.Unlock()

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
/*
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
*/
func (network *Network) createChannel(messageID *KademliaID, returnChannel chan interface{}) {
	network.mux.Lock()
	network.waitingAnswerList[*messageID] = returnChannel
	network.mux.Unlock()
}


func (network *Network) TimeoutWaiter(sleepTime int, returnChannel chan interface{}, messageID *KademliaID) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
	network.mux.Lock()
	network.waitingAnswerList[*messageID] = nil
	network.mux.Unlock()
	returnChannel <-false
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


