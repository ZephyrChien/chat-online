//+build !server

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	name           = flag.String("n", "guest", "you could reset it by /name")
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
		go netWriter(localConn, ch)
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
	go serverReader(conn)
	go netWriter(conn, ch)
	time.Sleep(1*time.Second)
	if *name != "guest" {
		ch <- "/name" + *name + "\n"
	}

	//listen client writer
	listener, err := net.Listen("tcp", *localHost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for {
		receiveLocalConn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go localReader(receiveLocalConn, ch)
	}
}

func serverReader(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		now := time.Now().Format("01-02 15:04:05")
		fmt.Printf("<%s>%s\n", now, input.Text())
	}
}

func localReader(conn net.Conn, ch chan<- string) {
	fmt.Println("connect to client writer", conn.RemoteAddr().String())
	input := bufio.NewScanner(conn)
	for input.Scan() {
		ch <- input.Text() + "\n"
	}
	fmt.Println("client writer", conn.RemoteAddr().String(), "close")
	conn.Close()
}

func netWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, msg)
	}
}
