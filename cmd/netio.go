package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
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

//InitListener up to if tls is enabled
func InitListener(network, address string, ssl bool) (l net.Listener) {
	if ssl {
		log.Print("listen with tls..")
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		tlslistener, err := tls.Listen("tcp", address, config)
		if err != nil {
			log.Fatal(err)
		}
		l = tlslistener
	} else {
		tcplistener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal(err)
		}
		l = tcplistener
	}
	return l
}

//InitConn up to if tls is enabled
func InitConn(network, address string, ssl bool) (conn net.Conn) {
	if ssl {
		fmt.Println("use tls")
		config := &tls.Config{
			InsecureSkipVerify: true,
		}
		tlsconn, err := tls.Dial("tcp", address, config)
		if err != nil {
			log.Fatal(err)
		}
		state := tlsconn.ConnectionState()
		var version string
		switch v := state.Version; v {
		case tls.VersionTLS13:
			version = "1.3"
		case tls.VersionTLS12:
			version = "1.2"
		case tls.VersionTLS11:
			version = "1.1"
		default:
			version = "1.0 or older"
		}
		fmt.Println("tls handshake succeed")
		fmt.Printf("tls version: %s\n", version)
		fmt.Println("=========================") // x25
		conn = tlsconn
	} else {
		tcpconn, err := net.Dial("tcp", address)
		if err != nil {
			log.Fatal(err)
		}
		conn = tcpconn
	}
	return conn
}
