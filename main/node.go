package main

import (
	//"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Node struct {
	data                   map[KademliaID]time.Time
	rt                     *RoutingTable
	mux                    *sync.Mutex
	republishTimeSeconds   int
	republishRandomSeconds int
}

// NewNode initializes node object. Setting republish time.
// Returns node as *Node
func NewNode(rt *RoutingTable) *Node {
	node := &Node{}
	node.rt = rt
	node.data = make(map[KademliaID]time.Time)
	node.mux = &sync.Mutex{}
	node.republishTimeSeconds = 15
	node.republishRandomSeconds = 10
	return node
}

// Store adds entry to node.data map.
// Adds entry as ID:timestamp. Timestamp represents when a republish is due.
func (node *Node) Store(key KademliaID, timeStamp time.Time) {
	node.mux.Lock()
	node.data[key] = (timeStamp.Add(time.Duration(node.republishTimeSeconds+rand.Intn(node.republishRandomSeconds)) * time.Second))
	node.mux.Unlock()
}

// GetDataMap returns map with data ID and timestamp.
// This is thread safe by using the mutex lock in node struct.
func (node *Node) GetDataMap() map[KademliaID]time.Time {
	node.mux.Lock()
	defer node.mux.Unlock()
	return node.data
}

// DeleteEntry deletes an entry from node.data map.
// Also removes the corresponding file in kademliastorage directory.
func (node *Node) DeleteEntry(dataEntryID KademliaID, storageMux *sync.Mutex) {
	node.mux.Lock()
	delete(node.data, dataEntryID)
	node.mux.Unlock()
	storageMux.Lock()
	err := os.Remove("kademliastorage/" + node.rt.me.ID.String() + "/" + dataEntryID.String())
	storageMux.Unlock()
	if err != nil {
		log.Fatal(err)
	}
}

// GotData checks node.data map if an entry exists.
// Returns true if it exists. False otherwise.
func (node *Node) GotData(key KademliaID) bool {
	node.mux.Lock()
	if _, ok := node.data[key]; ok {
		node.mux.Unlock()
		return true
	} else {
		node.mux.Unlock()
		return false
	}
}
