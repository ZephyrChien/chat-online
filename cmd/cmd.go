package cmd

import (
	"fmt"
	"log"
	"strings"
)

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
	case "list":
		for c := range clients {
			cli.Tunnel <- c.Name
		}
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
	switch val := ""; {
	case strings.HasPrefix(dat.Message, "/name"):
		fmt.Sscanf(dat.Message, "/name%s", &val)
		if val == "" {
			fmt.Println("/name: err args")
		} else {
			od.Is = true
			od.Key = "name"
			od.Val = val
			send()
		}
	case strings.HasPrefix(dat.Message, "/list"):
		fmt.Sscanf(dat.Message, "/name%s", &val)
		if val != "" {
			fmt.Println("/list: err args")
		} else {
			od.Is = true
			od.Key = "list"
			send()
		}
	default:
		fmt.Printf("unknown cmd: %s\n", dat.Message)
	}
}
