package main

import (
	"bufio"
	"io"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"path/filepath"
)

const BUFFERSIZE = 1024

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

func sendClientMode() {
	fmt.Println("Starting client send...")
	connection, error := net.Dial("tcp", "localhost:3000")
	if error != nil {
		fmt.Println(error)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func receiveClientMode() {
	// fmt.Println("Starting client recieve...")
	// connection, error := net.Dial("tcp", "localhost:3000")
	// if error != nil {
	// 	fmt.Println(error)
	// }
	// client := &Client{socket: connection}
	// go client.receive()
	go receiveFile()

	var input string
	fmt.Scanln(&input)

}

func receiveFile() {
	connection, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	path := filepath.Join("./file/", fileName)
	newFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}

func main() {
	flagMode := flag.String("mode", "send", "start client in send or receive mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "send" {
		sendClientMode()
	} else {
		receiveClientMode()
	}
}
