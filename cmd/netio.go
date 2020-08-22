package cmd

import (
	"fmt"
	"net"
	"os"
)

//RemoteWriter send data to partner
func RemoteWriter(conn net.Conn, ch <-chan string) {
	for str := range ch {
		fmt.Fprintln(conn, str)
	}
}

//RemoteEnWriter encrypt data then send to partner
func RemoteEnWriter(conn net.Conn, ch <-chan string, key string) {
	for msg := range ch {
		fmt.Fprintln(conn, Encrypt([]byte(msg), key))
	}
}

//PrintErr directly print err
func PrintErr(err error) {
	fmt.Fprintln(os.Stderr, err)
}

//ExDataWriter send client/server hello to the remote
func ExDataWriter(conn net.Conn, extra, key string) {
	dat := &Data{Extra: extra}
	crypt := Encrypt(MakeJSON(dat), key)
	fmt.Fprintln(conn, crypt)
}
