package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func newHash(str string) string {

	hashMachine := sha256.New()
	hashMachine.Write([]byte(str))
	hash_val := hex.EncodeToString(hashMachine.Sum(nil))
	return hash_val
}

func encrypter(fileid string, passphrase string, data []byte) []byte {
	key := fileid + passphrase
	newCipher, err := aes.NewCipher([]byte(newHash(key)))
	if err != nil {
		println("error when creating cipher", err)
		os.Exit(1)
	}
	counter, err := cipher.NewGCM(newCipher)
	nonce := make([]byte, counter.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	text := counter.Seal(nonce, nonce, data, nil)
	return text

}

func decrypter(data []byte, fileid string, passphrase string) []byte {
	decryptKey := []byte(newHash(fileid + passphrase))
	newCipher, _ := aes.NewCipher(decryptKey)
	counter, _ := cipher.NewGCM(newCipher)
	sizeOfNonce := counter.NonceSize()
	nonce, text := data[:sizeOfNonce], data[sizeOfNonce:]
	readable, _ := counter.Open(nil, nonce, text, nil)
	return readable
}
