//+build !server

package main

import (
	"bufio"
	//"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//arguments
var (
	username       = flag.String("n", "guest", "you could reset it by /name")
	localHost      = flag.String("l", "127.0.0.1:9000", "listen client writer")
	remoteHost     = flag.String("h", "127.0.0.1:8000", "remote server address")
	isClientWriter = flag.Bool("write", false, "start as client writer")
)

func init() {
	flag.Parse()
}
func main() {
	if *isClientWriter {
		ch := make(chan string, 1)
		localConn, err := net.Dial("tcp", *localHost)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer localConn.Close()
		go localWriter(localConn, ch)
		fmt.Println("start as client writer")
		inputReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("writer:|")
			input, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			ch <- input
		}
	}

	ch := make(chan []byte, 1)
	//connect to remote server
	conn, err := net.Dial("tcp", *remoteHost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()
	go remoteReader(conn)
	go remoteWriter(conn, ch)
	time.Sleep(1 * time.Second)

	if *username != "guest" {
		ch <- makeJSON(&data{CMD: command{true, "name", *username}})
	}

	//listen client writer
	listener, err := net.Listen("tcp", *localHost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		receiveLocalConn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go localReader(receiveLocalConn, ch)
	}
}

func localReader(conn net.Conn, ch chan<- []byte) {
	fmt.Println("connect to client writer", conn.RemoteAddr().String())
	input := bufio.NewScanner(conn)
	for input.Scan() {
		cmd := &command{}
		dat := &data{Name: *username, Message: input.Text()}
		if strings.HasPrefix(input.Text(), "/") {
			handleCMDC(cmd, dat)
		}
		buf := makeJSON(dat)
		ch <- buf
	}
	fmt.Println("client writer", conn.RemoteAddr().String(), "close")
	conn.Close()
}

func localWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, msg)
	}
}

func remoteReader(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		now := time.Now().Format("01-02 15:04:05")
		fmt.Printf("<%s>%s\n", now, input.Text())
	}
}

func remoteWriter(conn net.Conn, ch <-chan []byte) {
	for dat := range ch {
		conn.Write(dat)
	}
}

func handleCMDC(cmd *command, dat *data) {
	switch {
	case strings.HasPrefix(dat.Message, "/name"):
		fmt.Sscanf(dat.Message, "/name%s", &cmd.Val)
		if cmd.Val == "" {
			fmt.Println("/name: err args")
		} else {
			cmd.Is = true
			cmd.Key = "name"
		}
	default:
		fmt.Printf("unknown cmd: %s\n", dat.Message)
	}
	dat.Message = ""
	dat.CMD = *cmd
}
