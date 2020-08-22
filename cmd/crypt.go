package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

//Encrypt aes-128-gcm crypt, base64 encode
func Encrypt(dat []byte, key string) string {
	k := []byte(key)[:16]
	block, err := aes.NewCipher(k)
	if err != nil {
		PrintErr(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		PrintErr(err)
	}
	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		PrintErr(err)
	}
	crypted := aesgcm.Seal(nonce, nonce, dat, nil)
	return Base64Encode(crypted)
}

//Decrypt base64 decode, aes-128-gcm decypt
func Decrypt(encoded, key string) []byte {
	k := []byte(key)[:16]
	buf := Base64Decode(encoded)
	crypted, nonce := buf[12:], buf[:12]
	block, err := aes.NewCipher(k)
	if err != nil {
		PrintErr(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		PrintErr(err)
	}
	dat, err := aesgcm.Open(nil, nonce, crypted, nil)
	if err != nil {
		PrintErr(err)
	}
	return dat
}
