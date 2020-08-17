//+build server

package main

import (
	"flag"
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

var (
	port=flag.Int("p",8000,"source port")
	source=flag.String("s","0.0.0.0","source address")
	host string
)

func init(){
	flag.Parse()
	host=fmt.Sprintf("%s:%d",*source,*port)
}
func main() {
	listener, err := net.Listen("tcp", host)
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
			close(cli.tunnel)
			delete(clients, cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	src := conn.RemoteAddr().String()
	cli := client{ip: src, name: src, tunnel: ch}
	entering <- cli
	logWriter(cli.ip)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := input.Text()
		if strings.HasPrefix(msg, "/") {
			handleOrder(&cli, msg)
			continue
		}
		message <- fmt.Sprintf("[%s]  %s", cli.name, msg)
	}
	leaving <- cli
	message <- fmt.Sprintf("[world]  %s is dead!", cli.name)
	conn.Close()
}

func handleOrder(cli *client, msg string) {
	switch {
	case strings.HasPrefix(msg, "/name"):
		val := new(string)
		fmt.Sscanf(msg, "/name%s", val)
		if *val == "" {
			cli.tunnel <- "cmdErr: /name"
		} else {
			delete(clients, *cli)
			cli.name = *val
			entering <- *cli //sync map
			cli.tunnel <- "your name is " + cli.name
			message <- fmt.Sprintf("[world]  %s has arrived!", cli.name)
			logWriter(fmt.Sprintln(cli.ip,"=", cli.name))
		}
	default:
		cli.tunnel <- fmt.Sprintf("unknown cmd: %s", msg)
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
