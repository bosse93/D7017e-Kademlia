package main

import (
	"container/list"
	"fmt"
)

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushBack(contact)
		}
	} else {
		bucket.list.MoveToBack(element)
	}
}

func (bucket *bucket) AddContactNetwork(contact Contact, network *Network) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushBack(contact)
		} else {
			fmt.Println("Full bucket. Pinging")
			answerChannel := make(chan interface{})
			network.SendPingMessage(bucket.list.Remove(bucket.list.Front()).(Contact), answerChannel)
			select {
			case pingAnswer := <-answerChannel:
				switch pingAnswer := pingAnswer.(type) {
				case Contact:
					bucket.list.PushBack(pingAnswer)

				case bool:
					bucket.list.Remove(bucket.list.Front())
					bucket.list.PushBack(contact)
				}
			}
		}
	} else {
		bucket.list.MoveToBack(element)
	}
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
