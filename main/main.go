package main

import (
	"fmt"
	"strconv"
	"os"
	"D7024e-Kademlia/github.com/urfave/cli"
	//"time"
	"sort"
)

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name: "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "runs old main",
			Action:  func(c *cli.Context) error {
				runTest()
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				saved := ""
				if c.NArg() > 0 {
					saved = c.Args().Get(0)
				}
				if c.String("lang") == "spanish" {
					fmt.Println("hola", saved)
				} else {
					fmt.Println("Added", saved)
				}

				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
	
}

func runTest() {
	firstNode := NewContact(NewRandomKademliaID(), "localhost:8000")
	firstNodeRT := NewRoutingTable(firstNode)
	lastNetwork := NewNetwork(NewNode(firstNodeRT), "localhost", 8000)

	nodeList := []*RoutingTable{firstNodeRT}
	//lastNode := firstNode
	//create 100 nodes
	for i := 0; i < 100; i++ {
		port := 8001 + i
		a := "localhost:" + strconv.Itoa(port)


		ID := NewRandomKademliaID()
		rt := NewRoutingTable(NewContact(ID, a))
		nodeList = append(nodeList, rt)
		rt.AddContact(firstNodeRT.me)
		nw := NewNetwork(NewNode(rt), "localhost", port)
		fmt.Println("Ny Nod varv " + strconv.Itoa(i+1) + ": " + rt.me.String())
		//go nw.Listen("localhost", port)
		//time.Sleep(500 * time.Millisecond)
		kademlia := NewKademlia(nw)

		lookupResult := kademlia.LookupContact(ID, false)

		for q := range lookupResult {
			rt.AddContact(lookupResult[q])
		}
		lastNetwork = nw
	}

	printFirstNodeRT(firstNode, firstNodeRT)
	printLastNodeRT(nodeList)


	kademlia := NewKademlia(lastNetwork)
	done := kademlia.Store(NewKademliaID("FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "data to store")
	fmt.Println("done " + strconv.FormatBool(done))
	data := kademlia.LookupData("FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	fmt.Println("Data returned" + data)
	/*for k1, v := range IDRTList {
		for k2, v2 := range v.node.data {
			fmt.Println("Node " + k1.String() + " has " + v2 + " stored for key " + k2.String())
		}
	}*/

}

func printFirstNodeRT(firstNode Contact, firstNodeRT *RoutingTable) {
	fmt.Println("Node: " + firstNode.ID.String())
	for z := range firstNodeRT.buckets {
		contactList := firstNodeRT.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
}

func printLastNodeRT(nodeList []*RoutingTable) {
	lastNode := nodeList[len(nodeList)-1]
	fmt.Println("Node: " + lastNode.me.String())
	for z := range lastNode.buckets {
		contactList := lastNode.buckets[z]
		fmt.Println("Bucket: " + strconv.Itoa(z))
		for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
			contact := elt.Value.(Contact)
			fmt.Println(contact.String())
		}
	}
}

func printAllNodesRT(nodeList []*RoutingTable) {
	for w := range nodeList {
		fmt.Println("Node: " + nodeList[w].me.String())
		for z := range nodeList[w].buckets {
			contactList := nodeList[w].buckets[z]
			fmt.Println("Bucket: " + strconv.Itoa(z))
			for elt := contactList.list.Front(); elt != nil; elt = elt.Next() {
				contact := elt.Value.(Contact)
				fmt.Println(contact.String())
			}
		}
	}


}
