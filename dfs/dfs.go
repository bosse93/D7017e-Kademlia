package main

import (
  "os"
  "sort"
  "D7024e-Kademlia/github.com/urfave/cli"
  "fmt"
	"net"
	"bufio"
	protoPack "D7024e-Kademlia/proto"
	"log"
	"D7024e-Kademlia/github.com/protobuf/proto"
)
//FRONTEND CLI

//Make request to a node

func connect(Usage string, m *protoPack.WrapperMessage){
	p :=  make([]byte, 2048)

	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	sendPacket(marshalHelper(m), sourceAddress)
	//fmt.Fprintf(conn, Usage + arg0)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}

	conn.Close()
}

func main() {
  app := cli.NewApp()

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name:  "lang, l",
      Value: "english",
      Usage: "Language for the greeting",
    },
    cli.StringFlag{
      Name:  "config, c",
      Usage: "Load configuration from `FILE`",
    },
  }

  app.Commands = []cli.Command{
    {
      Name:    "store",
      Aliases: []string{"s", "Store", "S"},
      Usage:   "Store arg0 arg1",
      Action: func(c *cli.Context) error {
        if c.NArg() > 1 {
			packet := &protoPack.RequestStore{c.Args().Get(0), c.Args().Get(1)}
			wrapperMsg := &protoPack.WrapperMessage_RequestStore{packet}
			wrapper := &protoPack.WrapperMessage{"RequestContact", "frontend", "FFFFFFFF0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", wrapperMsg}

          //store c.Args().First()
          fmt.Println("Sending server request")
          connect("Store", wrapper)
        }
        return nil
      },
    },
    {
      Name:    "cat",
      Aliases: []string{"c"},
      Usage:   "Prints content of arg0",
      Action: func(c *cli.Context) error {
        if c.NArg() > 0 {
          //cat c.Args().First()


        }
        connect("Cat", "")
        return nil
      },
    },
    {
      Name:    "pin",
      Aliases: []string{"p"},
      Usage:   "Pins arg0",
      Action: func(c *cli.Context) error {
        if c.NArg() > 0 {
          //pin c.Args().First()

        }
        return nil
      },
    },
    {
      Name:    "unpin",
      Aliases: []string{"u"},
      Usage:   "Unpins arg0",
      Action: func(c *cli.Context) error {
        if c.NArg() > 0 {
          //unpin c.Args().First()

        }
        return nil
      },
    },
  }

  sort.Sort(cli.FlagsByName(app.Flags))
  sort.Sort(cli.CommandsByName(app.Commands))

  app.Run(os.Args)
}

func marshalHelper(wrapper *protoPack.WrapperMessage) []byte{
	data, err := protoPack.Marshal(wrapper)
	if err != nil {
		log.Fatal("Marshall Error: ", err)
	}
	return data
}

func sendPacket(data []byte, targetAddress *net.UDPAddr) {
	buf := []byte(data)
	_, err := listenConnection.WriteToUDP(buf, targetAddress)
	if err != nil {
		log.Println(err)
	}
}