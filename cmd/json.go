package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Command struct {
	Is  bool   `json:"apply"`
	Key string `json:"option"`
	Val string `json:"value"`
}
type Data struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	CMD     Command `json:"cmd"`
	Extra   string  `json:"extra"`
}

func MakeJSON(dat *Data) []byte {
	buf, err := json.Marshal(dat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return buf
}
func ResolvJSON(buf []byte,dat *Data){
	json.Unmarshal(buf,dat)
}

func Base64Encode(buf []byte)string{
	return base64.StdEncoding.EncodeToString(buf)
}

func Base64Decode(str string)[]byte{
	buf,err:=base64.StdEncoding.DecodeString(str)
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}
	return buf
}