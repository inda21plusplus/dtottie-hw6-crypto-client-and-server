package main

import (
	"net"
	"os"
	"strconv"
	"strings"
)

//this seems to be the broken function
func treeEval(path string) {
	path = strings.ReplaceAll(path, "\\", "/")
	pathArr := strings.Split(path, "/")
	pathArr = pathArr[:len(pathArr)-1]
	parentPath := strings.Join(pathArr[:], "/")

	contents, _ := rootFS.ReadDir(parentPath)
	newName := 0

	//add the filenames together
	for _, file := range contents {
		filename := strings.Split(file.Name(), "\\")
		filenumber, _ := strconv.Atoi(filename[len(filename)-1])
		newName += filenumber
	}
	pathArr[len(pathArr)-1] = strconv.Itoa(newName)
	newParentPath := strings.Join(pathArr[:], "/")
	rootFS.Rename(parentPath, newParentPath)
	if len(pathArr) > 1 {
		treeEval(newParentPath)
	}

}

//initialize filesystem at "init" command
func initializeFS() {
	rootFS.Mkdir("0", 0777)
	rootFS.Mkdir("0/1", 0777)
	rootFS.Mkdir("0/2", 0777)
	rootFS.Mkdir("0/1/1", 0777)
	rootFS.Mkdir("0/1/2", 0777)
	rootFS.Mkdir("0/1/1/3", 0777)
	f, err := rootFS.OpenFile("0/1/1/5", os.O_CREATE, 0)
	if err != nil {
		println("could not open file: ", err)
	}
	f.Close()

	treeEval("0/1/1/file")

	sessionPath = ""

}

//cd command
func changeDirectory(dir string) {
	sessionPath = sessionPath + "/" + dir
}

func writeFile(filename string, conn net.Conn) {

	f, err := rootFS.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		rootFS.OpenFile(filename, os.O_CREATE, 0)
		f, _ = rootFS.OpenFile(filename, os.O_RDWR, 0)
	}
	println(f)
	data := readInput(conn)
	f.Write([]byte(data))
}

func makeFile(filename string) error {
	filename = sessionPath + "/" + filename
	_, err := rootFS.OpenFile(filename, os.O_CREATE, 0)
	if err != nil {
		return err
	}
	return nil
}

func readFile(filename string, conn net.Conn) error {
	f, err := rootFS.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024)
	f.Read(buffer)
	conn.Write(buffer)
	return nil
}
