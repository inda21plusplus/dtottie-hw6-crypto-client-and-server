package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

const (
	RHOST = "127.0.0.1"
	RPORT = "8080"
	PASS  = "securepass123"
)

func startClient() {

	//initialize reader
	reader := bufio.NewReader(os.Stdin)

	hostname := RHOST + ":" + RPORT
	conn := connect(hostname)
	loopConnection(conn, reader)

}

func connect(hostname string) net.Conn {
	connAddr, err := net.ResolveTCPAddr("tcp", hostname)
	if err != nil {
		println(err)
		println("Invalid TCP address")
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, connAddr)
	if err != nil {
		println(err, "\nError connecting to the server")
	}
	println("connected to: ", hostname)
	return conn

}

func loopConnection(conn net.Conn, reader *bufio.Reader) {
	for {

		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			println(err)
		}
		println(string(msg))

		message, err := reader.ReadString('\n')
		message = message[:len(message)-2]

		if err != nil || message == "\n" || len(message) == 0 {
			conn.Write([]byte("test\n"))
		}

		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			println("failed to write")
		}

		//interpret arguments on client side
		arguments := strings.Split(message, " ")
		switch arguments[0] {
		case "quit":
			return
		case "read":
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				println(err)
			}
			println(interpret(msg, arguments[1]))
		case "write":
			message, err := reader.ReadString('\n')
			if err != nil {
				println(err)
			}
			conn.Write(encryptText(message, arguments[1]))
			conn.Write([]byte("\n"))
		}

	}
}

//decrypt received data
func interpret(data string, fileid string) string {
	text := decrypter(fileid, PASS, []byte(data))
	return string(text)
}

//encrypt data before sending
func encryptText(data string, fileid string) []byte {
	byteData := encrypter(fileid, PASS, []byte(data))
	return byteData
}
