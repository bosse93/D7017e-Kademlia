package main

import "fmt"

type Node struct {
	data map [KademliaID]string
	rt *RoutingTable
}

func NewNode(rt *RoutingTable) *Node  {
	node := &Node{}
	node.rt = rt
	node.data = make(map [KademliaID]string)
	return node
}

func (node *Node) Store(key KademliaID, data string)  {
	fmt.Println("Storing")
	node.data[key] = data
}
