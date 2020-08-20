package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type command struct {
	Is  bool   `json:"apply"`
	Key string `json:"option"`
	Val string `json:"value"`
}
type data struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	CMD     command `json:"cmd"`
	Extra   string  `json:"extra"`
}

func makeJSON(dat *data) []byte {
	buf, err := json.Marshal(dat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return buf
}
