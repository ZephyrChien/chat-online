package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//Client :specified by name
type Client struct {
	IP     string
	Name   string
	Extra  string
	Tunnel chan<- string
}

//Status :the combination of channel
type Status struct {
	Entering chan Client
	Leaving  chan Client
	Message  chan string
}

//server

//HandleCMDS handle command sent from client
func HandleCMDS(clients map[Client]bool, cli *Client, stat *Status, dat *Data, mlog *log.Logger) {
	switch key, val := dat.CMD.Key, dat.CMD.Val; key {
	case "name":
		delete(clients, *cli)
		cli.Name = val
		stat.Entering <- *cli //sync map
		cli.Tunnel <- "You are " + cli.Name
		stat.Message <- fmt.Sprintf("[world]%s has arrived!", cli.Name)
		mlog.Print(fmt.Sprintf("%s -> %s", cli.IP, cli.Name))
	default:
		cli.Tunnel <- fmt.Sprintf("unknown cmd: %s", key)
	}
}

//client

//HandleCMDC handle message like "/help"
//respond on client or send it to remote server
func HandleCMDC(dat *Data, od *Command, ch chan<- string, key string) {
	send := func() {
		dat.Message = ""
		dat.CMD = *od
		buf := MakeJSON(dat)
		ch <- Encrypt(buf, key)
	}
	switch {
	case strings.HasPrefix(dat.Message, "/name"):
		val := ""
		fmt.Sscanf(dat.Message, "/name%s", &val)
		if val == "" {
			fmt.Println("/name: err args")
		} else {
			od.Is = true
			od.Key = "name"
			od.Val = val
			send()
		}
	default:
		fmt.Printf("unknown cmd: %s\n", dat.Message)
	}
}

//io

//RemoteWriter send data to partner
func RemoteWriter(conn net.Conn, ch <-chan string) {
	for str := range ch {
		fmt.Fprintln(conn, str)
	}
}

//RemoteEnWriter encrypt data then send to partner
func RemoteEnWriter(conn net.Conn, ch <-chan string, key string) {
	for msg := range ch {
		fmt.Fprintln(conn, Encrypt([]byte(msg), key))
	}
}

//PrintErr directly print err
func PrintErr(err error) {
	fmt.Fprintln(os.Stderr, err)
}
