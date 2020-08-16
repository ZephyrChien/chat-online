//+build !server

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()
	go serverReader(conn)
	ch := make(chan string, 1)
	inputReader := bufio.NewReader(os.Stdin)
	go serverWriter(conn, ch)
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		ch <- input
		fmt.Printf("")
	}
}

func serverReader(conn net.Conn) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Printf("%s\n", input.Text())
	}
}

func serverWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, msg)
	}
}
