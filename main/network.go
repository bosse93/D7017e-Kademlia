package main

import (
	"D7024e-Kademlia/github.com/protobuf/proto"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Network struct {
	node                      *Node
	waitingAnswerList         map[KademliaID](chan interface{})
	listenConnection          *net.UDPConn
	mux                       *sync.Mutex
	rtMux                     *sync.Mutex
	timeoutTime               int
	fileNetwork               *FileNetwork
	republishSleepTimeSeconds int
}

// Initializes network object.
// NewNetwork sets a UDP listener on ip:port.
// Spawns go routine for RepublishData function.
// Returns network object.
func NewNetwork(node *Node, fileNetwork *FileNetwork, ip string, port int) *Network {
	network := &Network{}
	network.node = node
	network.waitingAnswerList = make(map[KademliaID]chan interface{})
	network.mux = &sync.Mutex{}
	network.rtMux = &sync.Mutex{}
	network.timeoutTime = 5
	network.fileNetwork = fileNetwork
	network.republishSleepTimeSeconds = 10

	// ESTABLISH UDP CONNECTION
	serverAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	network.listenConnection = serverConn
	buf := make([]byte, 4096)
	fmt.Println("Listening on port " + strconv.Itoa(port))
	go network.Listener(buf)
	go network.RepublishData()

	return network
}

// Listener waits for packets on listenConnection.
// Unmarshalls packet and spawns routine to handle the reply/request.
func (network *Network) Listener(buf []byte) {
	defer network.listenConnection.Close()

	//For each new packet do marshalling
	for {
		n, addr, _ := network.listenConnection.ReadFromUDP(buf)
		message := &WrapperMessage{}
		replyErr := proto.Unmarshal(buf[0:n], message)

		if message.ID[0:5] == "Reply" {
			go network.HandleReply(message, replyErr, addr)
		} else {
			go network.HandleRequest(message, replyErr, addr)
		}
	}
}

// SendPingMessage sends a RequestPing message to contact.
// Creates a RequestPing packet and marshalls it. Then sends it of to contact.
// Spawns new go routine to wait for answer.
func (network *Network) SendPingMessage(contact Contact, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestPing{}
	wrapperMsg := &WrapperMessage_RequestPing{packet}
	wrapper := &WrapperMessage{"RequestPing", network.node.rt.me.ID.String(), messageID.String(), wrapperMsg}

	network.CreateChannel(messageID, returnChannel)
	network.SendPacket(network.MarshalHelper(wrapper), remoteAddr)

	go network.TimeoutWaiter(network.timeoutTime, returnChannel, messageID)
}

// SendFindContactMessage sends a RequestContact message to contact.
// Creates a RequestContact packet and marshalls it. Then sends it of to contact.
// Spawns new go routine to wait for answer.
func (network *Network) SendFindContactMessage(targetID *KademliaID, contact *Contact, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestContact{targetID.String()}
	wrapperMsg := &WrapperMessage_RequestContact{packet}
	wrapper := &WrapperMessage{"RequestContact", network.node.rt.me.ID.String(), messageID.String(), wrapperMsg}

	network.CreateChannel(messageID, returnChannel)
	network.SendPacket(network.MarshalHelper(wrapper), remoteAddr)

	go network.TimeoutWaiter(network.timeoutTime, returnChannel, messageID)

}

// SendFindDataMessage sends a RequestData message to contact.
// Creates a RequestData packet and marshalls it. Then sends it of to contact.
// Spawns new go routine to wait for answer.
func (network *Network) SendFindDataMessage(hash string, contact *Contact, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

	packet := &RequestData{hash}
	wrapperMsg := &WrapperMessage_RequestData{packet}
	wrapper := &WrapperMessage{"RequestData", network.node.rt.me.ID.String(), messageID.String(), wrapperMsg}

	network.CreateChannel(messageID, returnChannel)
	network.SendPacket(network.MarshalHelper(wrapper), remoteAddr)

	go network.TimeoutWaiter(network.timeoutTime, returnChannel, messageID)
}

// SendStoreMessage sends a RequestStore message to address.
// Creates a RequestStore packet and marshalls it. Then sends it of to address.
// Spawns new go routine to wait for answer.
func (network *Network) SendStoreMessage(hash string, address string, returnChannel chan interface{}) {
	messageID := NewRandomKademliaID()
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	CheckError(err)
	data := "PLS_REMOVE_ME"

	packet := &RequestStore{hash, data} //EDIT ME
	wrapperMsg := &WrapperMessage_RequestStore{packet}
	wrapper := &WrapperMessage{"RequestStore", network.node.rt.me.ID.String(), messageID.String(), wrapperMsg}

	network.CreateChannel(messageID, returnChannel)
	network.SendPacket(network.MarshalHelper(wrapper), remoteAddr)

	go network.TimeoutWaiter(network.timeoutTime, returnChannel, messageID)
}

// HandleReply handles all message of type Reply.
// Reply types is ReplyPing, ReplyContactLost, ReplyData, ReplyStore.
// Check reply type and sends information through the corresponding channel.
func (network *Network) HandleReply(message *WrapperMessage, replyErr error, sourceAddress *net.UDPAddr) {
	if replyErr != nil {
		fmt.Println(message.ID)
		log.Println("Something went wrong in Listen, err: ", replyErr)
		return
	}

	answerChannel := network.GetAnswerChannel(NewKademliaID(message.RequestID))

	if answerChannel == nil {
		fmt.Println("Forged Reply or Timeout")
		return
	}

	switch message.ID {
	case "ReplyPing":
		contact := NewContact(NewKademliaID(message.GetReplyPing().GetID()), message.GetReplyPing().GetAddress())
		answerChannel <- contact
		return

	case "ReplyContactList":
		contactList := []Contact{}
		for i := range message.GetReplyContactList().GetContacts() {
			contactList = append(contactList, NewContact(NewKademliaID(message.GetReplyContactList().Contacts[i].GetID()), message.GetReplyContactList().Contacts[i].GetAddress()))
		}
		answerChannel <- contactList
		break

	case "ReplyData":
		answerChannel <- sourceAddress.String()
		break

	case "ReplyStore":
		answerChannel <- message.GetReplyStore().GetData()
		break

	default:
		fmt.Println("Not a valid Reply ID. ID: " + message.ID)
		return
	}
	go network.UpdateRoutingTable(message.SourceID, sourceAddress.String())

}

// HandleRequest handels all messages of type request.
// Request types is RequestPing, RequestContact, RequestData, RequestStore.
// Check request type and find corresponding data. Marshall it into Reply packet.
// Updates routingtable with source contact and send reply back.
func (network *Network) HandleRequest(message *WrapperMessage, replyErr error, sourceAddress *net.UDPAddr) {
	if replyErr != nil {
		fmt.Println(message.ID)
		log.Println("Something went wrong in Listen, err: ", replyErr)
		return
	}

	var wrapper *WrapperMessage

	switch message.ID {
	case "RequestPing":
		fmt.Println("Ping Recieved")
		packet := &ReplyPing{network.node.rt.me.ID.String(), network.node.rt.me.Address}
		wrapperMsg := &WrapperMessage_ReplyPing{packet}
		wrapper = &WrapperMessage{"ReplyPing", network.node.rt.me.ID.String(), message.RequestID, wrapperMsg}
		break

	case "RequestContact":
		network.rtMux.Lock()
		contactListReply := network.GetClosestContacts(message.GetRequestContact().GetTarget())
		network.rtMux.Unlock()

		packet := &ReplyContactList{contactListReply}
		wrapperMsg := &WrapperMessage_ReplyContactList{packet}
		wrapper = &WrapperMessage{"ReplyContactList", network.node.rt.me.ID.String(), message.RequestID, wrapperMsg}
		break

	case "RequestData":
		haveData := network.node.GotData(*NewKademliaID(message.GetRequestData().GetKey()))
		if haveData {
			packet := &ReplyData{"PLSREMOVEME"}
			wrapperMsg := &WrapperMessage_ReplyData{packet}
			wrapper = &WrapperMessage{"ReplyData", network.node.rt.me.ID.String(), message.RequestID, wrapperMsg}

		} else {
			network.rtMux.Lock()
			contactListReply := network.GetClosestContacts(message.GetRequestData().GetKey())
			network.rtMux.Unlock()

			packet := &ReplyContactList{contactListReply}
			wrapperMsg := &WrapperMessage_ReplyContactList{packet}
			wrapper = &WrapperMessage{"ReplyContactList", network.node.rt.me.ID.String(), message.RequestID, wrapperMsg}
		}

		break

	case "RequestStore":
		fileID := *NewKademliaID(message.GetRequestStore().GetKey())
		haveFile := network.node.GotData(fileID)
		if !haveFile {
			go network.fileNetwork.DownloadFile(&fileID, sourceAddress.String(), false)
		} else {
			network.node.Store(fileID, time.Now())
		}

		packet := &ReplyStore{"ok"}
		wrapperMsg := &WrapperMessage_ReplyStore{packet}
		wrapper = &WrapperMessage{"ReplyStore", network.node.rt.me.ID.String(), message.RequestID, wrapperMsg}
		break

	default:
		fmt.Println("Not a valid Request ID. ID: " + message.ID)
		return

	}
	go network.UpdateRoutingTable(message.SourceID, sourceAddress.String())
	network.SendPacket(network.MarshalHelper(wrapper), sourceAddress)
}

// CheckError prints err.
func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

// GetClosestContacts finds 20 closest contacts in routing table to target.
// Creates a ContactList packet for marshalling.
// Returns ContactList packet.
func (network *Network) GetClosestContacts(target string) (contactListReply []*ReplyContactList_Contact) {
	closestContacts := network.node.rt.FindClosestContacts(NewKademliaID(target), 20)
	contactListReply = []*ReplyContactList_Contact{}
	for i := range closestContacts {
		contactReply := &ReplyContactList_Contact{closestContacts[i].ID.String(), closestContacts[i].Address}
		contactListReply = append(contactListReply, contactReply)
	}
	return contactListReply
}

// CreateChannel adds returnChannel to waitingAnswerList.
func (network *Network) CreateChannel(messageID *KademliaID, returnChannel chan interface{}) {
	network.mux.Lock()
	network.waitingAnswerList[*messageID] = returnChannel
	network.mux.Unlock()
}

// GetAnswerChannel finds channel corresponding to requestID in waitingAnswerList.
// Returns channel if found in waitingAnswerList.
func (network *Network) GetAnswerChannel(requestID *KademliaID) (answerChannel chan interface{}) {
	network.mux.Lock()
	answerChannel = network.waitingAnswerList[*requestID]
	network.mux.Unlock()
	return answerChannel
}

// TimeoutWaiter will sleep for a timeout time then remove channel corresponding to messageID in waitingForAnswerList.
// Send false to channel indicating a timeout.
func (network *Network) TimeoutWaiter(sleepTime int, returnChannel chan interface{}, messageID *KademliaID) {
	time.Sleep(time.Duration(sleepTime) * time.Second)
	network.mux.Lock()
	network.waitingAnswerList[*messageID] = nil
	network.mux.Unlock()
	returnChannel <- false
}

// MarshalHelper marshalls wrapper packet.
// Returns marshalled packet.
func (network *Network) MarshalHelper(wrapper *WrapperMessage) []byte {
	data, err := proto.Marshal(wrapper)
	if err != nil {
		log.Fatal("Marshall Error: ", err)
	}
	return data
}

// SendPacket sends data packet to targetAddress through network.listenConnection connection.
func (network *Network) SendPacket(data []byte, targetAddress *net.UDPAddr) {
	buf := []byte(data)
	_, err := network.listenConnection.WriteToUDP(buf, targetAddress)
	if err != nil {
		log.Println(err)
	}
}

// UpdateRoutingTable creates a new Contact with contactID and contactAddress.
// Adds new Contact to RoutingTable
func (network *Network) UpdateRoutingTable(contactID string, contactAddress string) {
	contact := NewContact(NewKademliaID(contactID), contactAddress)
	network.rtMux.Lock()
	network.node.rt.AddContactNetwork(contact, network)
	network.rtMux.Unlock()
}

// RepublishData republish all data the node is responisble for to make sure data is replicated in the network.
// Sleeps for a time then check all entrys in nodes Data map.
// If entrys timestamp is less then current time a store message is sent to 20 closest node to data.
// If timestamp + time is less then current time entry is removed from Data map.
func (network *Network) RepublishData() {
	time.Sleep(time.Duration(network.republishSleepTimeSeconds) * time.Second)
	//fmt.Println("Republish Check")
	dataMap := network.node.GetDataMap()
	for dataEntryID, timestamp := range dataMap {
		if time.Now().After(timestamp) {
			kademlia := NewKademlia(network)
			contactList, _ := kademlia.LookupContact(&dataEntryID, false)
			i := 0

			for k := range contactList {
				fmt.Println("contact in republish" + DecodeHash(contactList[k].ID.String()))
				i++
				if contactList[k].ID.String() == network.node.rt.me.ID.String() {
					network.node.Store(dataEntryID, time.Now())
				} else {
					go network.SendStoreAndWaitForAnswer(dataEntryID.String(), contactList[k].Address, i)
				}
			}
			fmt.Println("Sent Republish")
		}
		if time.Now().After(timestamp.Add(time.Duration(2) * time.Second)) {
			go network.node.DeleteEntry(dataEntryID, network.fileNetwork.mux2)
		}
	}
	network.RepublishData()
}

// SendStoreAndWaitForAnswer sends a store message to address.
// Waits for confirmation on store message.
func (network *Network) SendStoreAndWaitForAnswer(dataEntryID string, address string, number int) {

	returnChannel := make(chan interface{})
	go network.SendStoreMessage(dataEntryID, address, returnChannel)
	returnValue := <-returnChannel
	switch returnValue := returnValue.(type) {
	case string:
		fmt.Println("Store " + strconv.Itoa(number) + " Reply: " + returnValue)
	case bool:
		fmt.Println("Store request timeout")
	default:
		fmt.Println("Something went wrong")
	}
}
