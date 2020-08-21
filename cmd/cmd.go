package cmd

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	IP     string
	Name   string
	Extra  string
	Tunnel chan<- string
}
type Status struct {
	Entering chan Client
	Leaving  chan Client
	Message  chan string
}


//server

func HandleCMDS(clients map[Client]bool, cli *Client, stat *Status, dat *Data) {
	switch key,val:=dat.CMD.Key,dat.CMD.Val;key{
	case "name":
		delete(clients, *cli)
		cli.Name = val
		stat.Entering <- *cli //sync map
		cli.Tunnel <- "your name is " + cli.Name
		stat.Message <- fmt.Sprintf("[world]%s has arrived!", cli.Name)
		LogWriter(fmt.Sprintln(cli.IP, "=", cli.Name))
	default:
		cli.Tunnel <- fmt.Sprintf("unknown cmd: %s", val)
	}
}
func LogWriter(msg string) {
	log.Printf("%s", msg)
}

//client
func HandleCMDC(dat *Data, od *Command) {
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
		}
	default:
		fmt.Printf("unknown cmd: %s\n", dat.Message)
	}
	dat.Message = ""
	dat.CMD = *od
}
//io
func RemoteWriter(conn net.Conn, ch <-chan string) {
	for str := range ch {
		fmt.Fprintln(conn, str)
	}
}