//+build !server

package main

import (
	"./cmd"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	ssl            = flag.Bool("ssl", false, "enable tls")
	key            = flag.String("k", "0000000000000000", "crypt key")
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
			cmd.PrintErr(err)
		}
		defer localConn.Close()
		go localWriter(localConn, ch)
		fmt.Println("start as client writer")
		inputReader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("writer:|")
			input, err := inputReader.ReadString('\n')
			if err != nil {
				cmd.PrintErr(err)
				continue
			}
			ch <- input
		}
	}
	ch := make(chan string, 1)
	//connect to remote server
	conn := cmd.InitConn("tcp", *remoteHost, *ssl)
	defer conn.Close()
	fmt.Println("client init")
	time.Sleep(1 * time.Second)
	fmt.Println("send client hello")
	cmd.ExDataWriter(conn, "c_hello", *key)
	fmt.Println("wait ack from server")
	if !cmd.OneTouchAuth(conn, "s_ack", *key, 2, log.New(os.Stderr, "", log.LstdFlags)) {
		log.Fatal("Auth failed")
	}
	fmt.Println("Auth successfully!")
	fmt.Println("=========================") // x25

	go remoteReader(conn)
	go cmd.RemoteWriter(conn, ch)
	time.Sleep(1 * time.Second)
	if *username != "guest" {
		dat := &cmd.Data{CMD: cmd.Command{Is: true, Key: "name", Val: *username}}
		ch <- cmd.Encrypt(cmd.MakeJSON(dat), *key)
	}

	//listen client writer
	listener, err := net.Listen("tcp", *localHost)
	if err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		receiveLocalConn, err := listener.Accept()
		if err != nil {
			cmd.PrintErr(err)
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
			cmd.HandleCMDC(dat, od, ch, *key)
			continue
		}
		buf := cmd.MakeJSON(dat)
		ch <- cmd.Encrypt(buf, *key)
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
		msg := cmd.Decrypt(input.Text(), *key)
		now := time.Now().Format("01-02 15:04:05")
		fmt.Printf("[%s]%s\n", now, msg)
	}
}
