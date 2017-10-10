package main

import (
	//"fmt"
	"time"
	"sync"
	"os"
	"log"
	"math/rand"
)

type Node struct {
	data map [KademliaID]time.Time
	rt *RoutingTable
	mux *sync.Mutex
	republishTimeSeconds int
	republishRandomSeconds int
}

func NewNode(rt *RoutingTable) *Node  {
	node := &Node{}
	node.rt = rt
	node.data = make(map [KademliaID]time.Time)
	node.mux = &sync.Mutex{}
	node.republishTimeSeconds = 15
	node.republishRandomSeconds = 10
	return node
}

func (node *Node) Store(key KademliaID, timeStamp time.Time) {
	node.mux.Lock()
	node.data[key] = (timeStamp.Add(time.Duration(node.republishTimeSeconds+ rand.Intn(node.republishRandomSeconds)) * time.Second))
	node.mux.Unlock()
}

func (node *Node) getDataMap() map [KademliaID]time.Time{
	node.mux.Lock()
	defer node.mux.Unlock()
	return node.data
}

func (node *Node) deleteEntry(dataEntryID KademliaID, storageMux *sync.Mutex) {
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

func (node *Node) gotData(key KademliaID) bool{
	node.mux.Lock()
	if _, ok := node.data[key]; ok{
		node.mux.Unlock()
		return true
	} else {
		node.mux.Unlock()
		return false
	}
}

