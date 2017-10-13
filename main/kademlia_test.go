package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestKademlia_FindNextNodeToAsk(t *testing.T) {
	// GIVEN
	localKademlia := NewKademlia(network)
	contacts := []Contact{}
	contacts = append(contacts, NewContact(NewRandomKademliaID(), "localhost:8200"))
	localKademlia.closest = NewContactCandidates()
	localKademlia.closest.Append(contacts)
	localKademlia.asked[*contacts[0].ID] = true
	// WHEN
	contact, success := localKademlia.FindNextNodeToAsk()
	// THEN
	if contact != nil {
		t.Error("Expected nil, got ", contact)
	}
	if success != false {
		t.Error("Expected false, got ", success)
	}

	localKademlia.asked[*contacts[0].ID] = false
	// WHEN
	contact2, success2 := localKademlia.FindNextNodeToAsk()
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
	kademlia := NewKademlia(network)
	contacts, _ := kademlia.LookupContact(HashKademliaID("100"), false)
	distances := NewContactCandidates()
	for i := 0; i < 100; i++ {
		dist := HashKademliaID(strconv.Itoa(i)).CalcDistance(HashKademliaID("100"))
		cont := []Contact{}
		c := NewContact(dist, "")
		c.distance = dist
		cont = append(cont, c)
		distances.Append(cont)
	}
	distances.Sort()

	for i := range contacts {
		if contacts[i].ID.CalcDistance(HashKademliaID("100")).String() != distances.contacts[i].ID.String() {
			t.Error("Expected distance to be " + distances.contacts[i].ID.String() + ", got " + contacts[i].ID.CalcDistance(HashKademliaID("100")).String())
		}
	}
}

func TestKademlia_LookupData(t *testing.T) {
	kademlia := NewKademlia(network)
	found := kademlia.LookupData("100")

	if found {
		t.Error("Expected found to be false, got" + strconv.FormatBool(found))
	}
}

func TestKademlia_Store(t *testing.T) {
	kademlia := NewKademlia(network)
	upload(network.node.rt.me.ID.String(), "testStore.txt")
	go kademlia.Store("testStore.txt")
	time.Sleep(time.Duration(1) * time.Second)
	kademlia2 := NewKademlia(network)
	found := kademlia2.LookupData("testStore.txt")
	fmt.Println("found " + strconv.FormatBool(found))
	if !found {
		t.Error("Expected found to be 'true', got " + strconv.FormatBool(found))
	}
}
