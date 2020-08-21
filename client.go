//+build !server

package main

import (
	"./cmd"
	"bufio"
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

	ch := make(chan string, 1)
	//connect to remote server
	conn, err := net.Dial("tcp", *remoteHost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()
	go remoteReader(conn)
	go cmd.RemoteWriter(conn, ch)
	time.Sleep(1 * time.Second)

	if *username != "guest" {
		dat := &cmd.Data{CMD: cmd.Command{true, "name", *username}}
		ch <- cmd.Base64Encode(cmd.MakeJSON(dat))
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

func localReader(conn net.Conn, ch chan<- string) {
	fmt.Println("connect to client writer", conn.RemoteAddr().String())
	input := bufio.NewScanner(conn)
	for input.Scan() {
		od := new(cmd.Command)
		dat := &cmd.Data{Name: *username, Message: input.Text()}
		if strings.HasPrefix(input.Text(), "/") {
			cmd.HandleCMDC(dat, od)
		}
		buf := cmd.MakeJSON(dat)
		ch <- cmd.Base64Encode(buf)
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
