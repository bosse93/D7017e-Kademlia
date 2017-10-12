package main

import (
	"fmt"
	"testing"
)

var rt *RoutingTable

func TestRoutingTable(t *testing.T) {
	rt = NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8003"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8004"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8005"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8006"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8007"))

	//saknade 7% coverage, detta löser allt #trial&error (rad 37 routing table, 		if bucketIndex-i >= 0 {
	rt.FindClosestContacts(NewKademliaID("f000000000000000000000000000000000000000"), 20)

	contacts2 := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts2 {
		fmt.Println(contacts2[i].String())
		//t.Log(contacts2[i].String())
	}

}

func TestBucket_Len(t *testing.T) {
	burk := newBucket()
	burk.Len()
}

func TestRoutingTable_FindClosestContacts(t *testing.T) {
	rt.FindClosestContacts(NewKademliaID("FFFFFF5F00000000000000000000000000000000"), 5)
}
