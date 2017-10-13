package main

import (
	"D7024e-Kademlia/github.com/urfave/cli"
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	//protoPack "D7024e-Kademlia/proto"
	//"log"
	//"D7024e-Kademlia/github.com/protobuf/proto"
	"bytes"
	"io"
	"net/http"
	"strings"
)

//FRONTEND CLI

//Make request to a node

func connect(m string) {
	p := make([]byte, 2048)
	split := strings.Split(m, " ")
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, m)
	if split[0] == "store" {
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	} else if split[0] == "cat" {
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			n := bytes.IndexByte(p, 0)
			/*dat, err := ioutil.ReadFile(string(p[:n]))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("file content: " + string(dat))*/
			fmt.Println("file content: " + string(p[:n]))
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	} else if split[0] == "pin" {

	} else if split[0] == "unpin" {

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
			Usage:   "Store file",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					m := "store " + c.Args().Get(0)
					//store c.Args().First()
					connect(m)
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
					m := "cat " + c.Args().Get(0)
					//store c.Args().First()
					connect(m)
				}
				return nil
			},
		},
		{
			Name:    "pin",
			Aliases: []string{"p"},
			Usage:   "Pins arg0",
			Action: func(c *cli.Context) error {
				if c.NArg() > 0 {
					m := "pin " + c.Args().Get(0)
					connect(m)
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
					m := "unpin " + c.Args().Get(0)
					connect(m)
				}
				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"Start"},
			Usage:   "start network",
			Action: func(c *cli.Context) error {
				connect("start")
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func downloadFile(filepath string, url string) (err error) {
	fmt.Println("filepath: " + filepath + " url: " + url)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
