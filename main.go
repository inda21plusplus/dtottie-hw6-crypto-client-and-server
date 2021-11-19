package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("1 for Server\n2 for Client\n")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-2]

	switch text {
	case "1":
		startServer()
	case "2":
		startClient()
	}
}
