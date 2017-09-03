package main

import (
		"fmt"
		sample "D7017e-Kademlia/SampleCode"
)

func main() {
    fmt.Printf("hello, world\n")
    rt := sample.NewContact(sample.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
    fmt.Printf(rt.String())
}