//+build server

package main

import (
	"./cmd"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

// arguments
var (
	host   string
	port   = flag.Int("p", 8000, "source port")
	source = flag.String("s", "0.0.0.0", "source address")
)

var (
	stat    = new(cmd.Status)
	clients = make(map[cmd.Client]bool)
)

func init() {
	flag.Parse()
	host = fmt.Sprintf("%s:%d", *source, *port)
	stat.Entering = make(chan cmd.Client, 1)
	stat.Leaving = make(chan cmd.Client, 1)
	stat.Message = make(chan string, 1)
}
func main() {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcast() {
	for {
		select {
		case msg := <-stat.Message:
			cmd.LogWriter(msg)
			for cli := range clients {
				cli.Tunnel <- msg
			}
		case cli := <-stat.Entering:
			clients[cli] = true
		case cli := <-stat.Leaving:
			close(cli.Tunnel)
			delete(clients, cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string, 1)
	src := conn.RemoteAddr().String()
	cli := cmd.Client{IP: src, Name: src, Tunnel: ch}
	go cmd.RemoteWriter(conn, ch)
	stat.Entering <- cli
	cmd.LogWriter(cli.IP)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		dat := new(cmd.Data)
		cmd.ResolvJSON(cmd.Base64Decode(input.Text()), dat)
		if dat.CMD.Is {
			cmd.HandleCMDS(clients, &cli, stat, dat)
			continue
		}
		stat.Message <- fmt.Sprintf("[%s]|%s", cli.Name, dat.Message)
	}
	stat.Leaving <- cli
	stat.Message <- fmt.Sprintf("[world]%s is dead!", cli.Name)
	conn.Close()
}
