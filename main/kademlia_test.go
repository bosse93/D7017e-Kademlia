package main

import (
	"testing"
)

var kademlia = NewKademlia(network)

func TestKademlia_findNextNodeToAsk(t *testing.T) {
	// GIVEN
	localKademlia := kademlia
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewRandomKademliaID(), "localhost:8200"))
	localKademlia.closest = NewContactCandidates()
	localKademlia.closest.Append(contacts)
	localKademlia.asked[*contacts[0].ID] = true
	// WHEN
	contact, success := localKademlia.findNextNodeToAsk()
	// THEN
	if contact != nil {
		t.Error("Expected nil, got ", contact)
	}
	if success != false {
		t.Error("Expected false, got ", success)
	}

	localKademlia.asked[*contacts[0].ID] = false
	// WHEN
	contact2, success2 := localKademlia.findNextNodeToAsk()
	// THEN
	if contact2.ID != contacts[0].ID {
		t.Error("Expected " + contacts[0].String() + ", got " + contact.String())
	}
	if success2 != true {
		t.Error("Expected true, got ", success)
	}
}

func TestKademlia_AskNextNode(t *testing.T) {

}

func TestKademlia_LookupContact(t *testing.T) {
}

func TestKademlia_LookupData(t *testing.T) {

}

func TestKademlia_Store(t *testing.T) {

}
