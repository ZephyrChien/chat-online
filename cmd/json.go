package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

//Command :empty unless the message starts with "/"
type Command struct {
	Is  bool   `json:"apply"`
	Key string `json:"option"`
	Val string `json:"value"`
}

//Data :client send data in this format
type Data struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	CMD     Command `json:"cmd"`
	Extra   string  `json:"extra"`
}

//MakeJSON format key:val pairs as json
func MakeJSON(dat *Data) []byte {
	buf, err := json.Marshal(dat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return buf
}

//ResolvJSON read json into a variable
func ResolvJSON(buf []byte, dat *Data) {
	json.Unmarshal(buf, dat)
}

//Base64Encode encode bytes to string
func Base64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

//Base64Decode decode base64 string to bytes
func Base64Decode(str string) []byte {
	buf, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return buf
}
