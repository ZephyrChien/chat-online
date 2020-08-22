//+build server

package main

import (
	"./cmd"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// arguments
var (
	host    string
	port    = flag.Int("p", 8000, "source port")
	key     = flag.String("k", "0000000000000000", "crypt key")
	source  = flag.String("s", "0.0.0.0", "source address")
	logfile = flag.String("log", "access.log", "file to store log")
)

var (
	stat    = new(cmd.Status)
	clients = make(map[cmd.Client]bool)
	mlog    *log.Logger
	logFile *os.File
)

func init() {
	//get args
	flag.Parse()
	host = fmt.Sprintf("%s:%d", *source, *port)
	stat.Entering = make(chan cmd.Client, 1)
	stat.Leaving = make(chan cmd.Client, 1)
	stat.Message = make(chan string, 1)

	//init logger
	logFile, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	mlog = log.New(io.MultiWriter(os.Stderr, logFile), "", log.LstdFlags)
}
func main() {
	defer logFile.Close()
	listener, err := net.Listen("tcp", host)
	if err != nil {
		mlog.Fatal(err)
	}
	defer listener.Close()
	go broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			mlog.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcast() {
	for {
		select {
		case msg := <-stat.Message:
			mlog.Print(msg)
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
	go cmd.RemoteEnWriter(conn, ch, *key)
	stat.Entering <- cli
	mlog.Print(cli.IP)
	if cmd.OneTouchAuth(conn, "c_hello", *key, 3, mlog) {
		cmd.ExDataWriter(conn, "s_ack", *key)
		input := bufio.NewScanner(conn)
		for input.Scan() {
			dat := new(cmd.Data)
			plaintext, err := cmd.ServerDecrypt(input.Text(), *key, mlog)
			if err != nil {
				break
			}
			cmd.ResolvJSON(plaintext, dat)
			if dat.CMD.Is {
				cmd.HandleCMDS(clients, &cli, stat, dat, mlog)
				continue
			}
			stat.Message <- fmt.Sprintf("[%s]|%s", cli.Name, dat.Message)
		}
	}
	stat.Leaving <- cli
	stat.Message <- fmt.Sprintf("[world]%s is dead!", cli.Name)
	conn.Close()
}
