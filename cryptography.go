package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
)

func newHash(str string) string {

	hashMachine := md5.New()
	hashMachine.Write([]byte(str))
	hash_val := hex.EncodeToString(hashMachine.Sum(nil))
	return hash_val
}

func encrypter(fileid string, passphrase string, data []byte) []byte {
	key := fileid + passphrase
	keyhash := newHash(key)
	println("hash is: ", keyhash)
	newCipher, err := aes.NewCipher([]byte(keyhash))
	println(newCipher)
	if err != nil {
		println("error when creating cipher", err)
		os.Exit(1)
	}
	counter, err := cipher.NewGCM(newCipher)
	if err != nil {
		println("error when creating GCM", err)
		os.Exit(1)
	}
	nonce := make([]byte, counter.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	text := counter.Seal(nonce, nonce, data, nil)
	return text

}

func decrypter(fileid string, passphrase string, data []byte) []byte {
	decryptKey := []byte(newHash(fileid + passphrase))
	newCipher, _ := aes.NewCipher(decryptKey)
	counter, _ := cipher.NewGCM(newCipher)
	sizeOfNonce := counter.NonceSize()
	nonce, text := data[:sizeOfNonce], data[sizeOfNonce:]
	readable, _ := counter.Open(nil, nonce, text, nil)
	return readable
}
