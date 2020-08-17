//+build server

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	ip     string
	name   string
	extra  string
	tunnel chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	message  = make(chan string)
	clients  = make(map[client]bool)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
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
		case msg := <-message:
			logWriter(msg)
			for cli := range clients {
				cli.tunnel <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.tunnel)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	src := conn.RemoteAddr().String()
	cli := client{ip: src, name: src, tunnel: ch}
	entering <- cli
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := input.Text()
		if strings.HasPrefix(msg, "/") {
			handleOrder(&cli, msg)
			continue
		}
		message <- fmt.Sprintf("[%s]:  %s", cli.name, msg)
	}
	leaving <- cli
	message <- cli.name + " is dead!"
	conn.Close()
}

func handleOrder(cli *client, msg string) {
	if strings.HasPrefix(msg, "/name") {
		val := new(string)
		fmt.Sscanf(msg, "/name%s", val)
		if *val == "" {
			cli.tunnel <- "cmdErr: /name"
		} else {
			cli.name = *val
			cli.tunnel <- "your name is " + cli.name
			message <- cli.name + " has arrived!"
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
func logWriter(msg string) {
	log.Printf("%s", msg)
}
