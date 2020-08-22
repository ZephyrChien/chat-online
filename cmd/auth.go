package cmd

import (
	"bufio"
	"log"
	"net"
	"time"
)

//OneTouchAuth receive first packet from remote partner within a certain second
func OneTouchAuth(conn net.Conn, passwd, key string, timeout int, mlog *log.Logger) bool {
	ch := make(chan int)
	go func() {
		dat := new(Data)
		input := bufio.NewScanner(conn)
		input.Scan()
		if input.Text() == "" {
			return
		}
		plaintext, err := ServerDecrypt(input.Text(), key, mlog)
		if err != nil {
			return
		}
		ResolvJSON(plaintext, dat)
		if dat.Extra == passwd {
			ch <- -1
		}

	}()
	select {
	case <-ch:
		return true
	case <-time.After(time.Duration(timeout) * time.Second):
		return false
	}
}
