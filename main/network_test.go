package main

import (
	"testing"
	"net"
)

func TestNetwork_SendPingMessage(t *testing.T) {

}

func TestNetwork_HandleReplyPing(t *testing.T) {
	id := NewKademliaID("ffffffff00000000000000000000000000000000")
	channel := make(chan interface{})
	network.createChannel(id, channel)
	contact := NewContact(id, "localhost:8000")

	packet := &ReplyPing{contact.ID.String(), contact.Address}
	wrapperMsg := &WrapperMessage_ReplyPing{packet}
	wrapper := &WrapperMessage{"ReplyPing", id.String(), id.String(), wrapperMsg}
	addr, _ := net.ResolveUDPAddr("udp", contact.Address)

	go network.HandleReply(wrapper, nil, addr)

	x := <- network.getAnswerChannel(id)
	returnedContact, ok := x.(Contact)
	if ok {
		if returnedContact.ID.String() != "ffffffff00000000000000000000000000000000" {
			t.Error("Expected ffffffff00000000000000000000000000000000, got ", returnedContact.ID.String())
		}
		if returnedContact.Address != "localhost:8000" {
			t.Error("Expected localhost:8000, got " + returnedContact.Address)
		}
	} else {
		t.Error("Expected return to be of type 'Contact'")
	}
}

func TestNetwork_HandleReplyContactList(t *testing.T) {
	id := NewKademliaID("ffffffff00000000000000000000000000000000")
	channel := make(chan interface{})
	network.createChannel(id, channel)
	contact := NewContact(id, "localhost:8000")

	contactReply := &ReplyContactList_Contact{contact.ID.String(), contact.Address}
	contactListReply := []*ReplyContactList_Contact{}
	contactListReply = append(contactListReply, contactReply)

	packet := &ReplyContactList{contactListReply}
	wrapperMsg := &WrapperMessage_ReplyContactList{packet}
	wrapper := &WrapperMessage{"ReplyContactList", id.String(), id.String(), wrapperMsg}

	addr, _ := net.ResolveUDPAddr("udp", contact.Address)

	go network.HandleReply(wrapper, nil, addr)

	x := <- network.getAnswerChannel(id)
	returnedContacts, ok := x.([]Contact)
	if ok {
		if returnedContacts[0].ID.String() != "ffffffff00000000000000000000000000000000" {
			t.Error("Expected ffffffff00000000000000000000000000000000, got ", returnedContacts[0].ID.String())
		}
		if returnedContacts[0].Address != "localhost:8000" {
			t.Error("Expected localhost:8000, got " + returnedContacts[0].Address)
		}
	} else {
		t.Error("Expected return to be of type '[]Contact'")
	}
}

func TestNetwork_HandleReplyData(t *testing.T) {
	id := NewKademliaID("ffffffff00000000000000000000000000000000")
	channel := make(chan interface{})
	network.createChannel(id, channel)
	contact := NewContact(id, "localhost:8000")

	packet := &ReplyData{""}
	wrapperMsg := &WrapperMessage_ReplyData{packet}
	wrapper := &WrapperMessage{"ReplyData", id.String(), id.String(), wrapperMsg}

	addr, _ := net.ResolveUDPAddr("udp", contact.Address)

	go network.HandleReply(wrapper, nil, addr)

	x := <- network.getAnswerChannel(id)
	returnedAdress, ok := x.(string)
	if ok {
		if returnedAdress != "127.0.0.1:8000" {
			t.Error("Expected 127.0.0.1:8000, got ", returnedAdress)
		}
	} else {
		t.Error("Expected return to be of type 'String'")
	}
}

func TestNetwork_HandleReplyStore(t *testing.T) {
	id := NewKademliaID("ffffffff00000000000000000000000000000000")
	channel := make(chan interface{})
	network.createChannel(id, channel)
	contact := NewContact(id, "localhost:8000")

	packet := &ReplyStore{"ok"}
	wrapperMsg := &WrapperMessage_ReplyStore{packet}
	wrapper := &WrapperMessage{"ReplyStore", id.String(), id.String(), wrapperMsg}

	addr, _ := net.ResolveUDPAddr("udp", contact.Address)

	go network.HandleReply(wrapper, nil, addr)

	x := <- network.getAnswerChannel(id)
	reply, ok := x.(string)
	if ok {
		if reply != "ok" {
			t.Error("Expected ok, got ", reply)
		}
	} else {
		t.Error("Expected return to be of type 'String'")
	}
}

func TestNetwork_RepublishData(t *testing.T) {

}

func TestNetwork_HandleRequest(t *testing.T) {

}

func TestNetwork_SendFindContactMessage(t *testing.T) {

}

func TestNetwork_SendFindDataMessage(t *testing.T) {

}

func TestNetwork_SendStoreMessage(t *testing.T) {

}