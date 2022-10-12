package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

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
	fmt.Println("Starting client recieve...")
	connection, error := net.Dial("tcp", "localhost:3000")
	if error != nil {
		fmt.Println(error)
	}
	client := &Client{socket: connection}
	go client.receive()

	var input string
	fmt.Scanln(&input)

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
