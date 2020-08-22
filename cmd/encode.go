package cmd

import (
	"encoding/base64"
	"encoding/json"
)

//MakeJSON format key:val pairs as json
func MakeJSON(dat *Data) []byte {
	buf, err := json.Marshal(dat)
	if err != nil {
		PrintErr(err)
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
		panic(err)
	}
	return buf
}

//ServerBase64Decode close conn if unable to decode
func ServerBase64Decode(str string) (buf []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buf, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return buf, err
}
