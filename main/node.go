package main

import (
	"fmt"
	"time"
)

type Node struct {
	data map [KademliaID]string
	dataRepublishTime map [KademliaID]time.Time
	rt *RoutingTable
}

func NewNode(rt *RoutingTable) *Node  {
	node := &Node{}
	node.rt = rt
	node.data = make(map [KademliaID]string)
	node.dataRepublishTime = make(map [KademliaID]time.Time)
	return node
}

func (node *Node) Store(key KademliaID, data string)  (haveFile bool){
	if _, ok := node.data[key]; ok{
		fmt.Println("Updated refresh timer")
		node.dataRepublishTime[key] = (time.Now().Add(time.Duration(20) * time.Second))
		haveFile = true
		return haveFile
	} else {
		haveFile = false
		return haveFile
		/*
		node.data[key] = data
		node.dataRepublishTime[key] = (time.Now().Add(time.Duration(20) * time.Second))
		*/
	}
	
}