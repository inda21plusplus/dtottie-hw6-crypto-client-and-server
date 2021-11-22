package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/blang/vfs/memfs"
)

const (
	LHOST = "localhost"
	LPORT = "8080"
	TYPE  = "tcp"
)

//declare rootFS as global variable as this should be persistent and accessible all over the function

var sessionPath string
var rootFS *memfs.MemFS

//starts the TCP server
func startServer() {
	//set up in memory file system
	rootFS = memfs.Create()

	//open listener and verify that its running
	listener, err := net.Listen(TYPE, LHOST+":"+LPORT)
	if err != nil {
		fmt.Println("error when setting up listener")
		os.Exit(1)
	}
	//close listener when app closes
	defer listener.Close()
	fmt.Println("listening on " + LHOST + ":" + LPORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection ", err.Error())
			os.Exit(1)
		}
		println("connection received")
		conn.Write([]byte("connected\n"))
		handleRequest(conn)
	}
}

//handles the incoming request and listens for commands.
func handleRequest(conn net.Conn) {
	for {

		message := readInput(conn)

		arguments := strings.Split(message, " ")
		println(arguments[0])
		arguments[0] = strings.ReplaceAll(arguments[0], "\n", "")
		switch arguments[0] {
		case "init":
			initializeFS()
			conn.Write([]byte("OK\n"))
		case "write":
			writeFile(arguments[1], conn)
		case "read":
			filename := sessionPath + "/" + arguments[1]
			println(filename)
			readFile(filename, conn)
			println("returned from func")
		case "touch":
			makeFile(arguments[1])
		case "quit":
			println("Closed connection: ", conn.RemoteAddr().String())
			conn.Close()
			return
		case "pwd":
			conn.Write([]byte(sessionPath + "\n"))
		case "ls":
			println("current path: ", sessionPath)
			readPath, err := rootFS.ReadDir(sessionPath)
			if err != nil {
				println(err)
			}
			for _, element := range readPath {
				println("folder:", element.Name())
				conn.Write([]byte(element.Name() + "  "))

			}
			conn.Write([]byte("\n"))
		case "cd":
			arguments[1] = strings.ReplaceAll(arguments[1], "\n", "")
			switch arguments[1] {
			case "..":
				sessionPathTMP := strings.Split(sessionPath, "/")
				sessionPathTMP = sessionPathTMP[:len(sessionPathTMP)-1]
				sessionPath = strings.Join(sessionPathTMP, "/")
				conn.Write([]byte("new path: " + sessionPath + "\n"))

			default:
				changeDirectory(arguments[1])
				conn.Write([]byte("new path: " + sessionPath + "\n"))
			}

		default:
			println("unknown command")
			conn.Write([]byte("unknown command\n"))
		}
	}
}

func readInput(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		println("failed to read message from client", err)
		return ""
	}
	return message
}
