package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const BUFFERSIZE = 1024

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type ClientManager1 struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

func sendFileToClient(connection net.Conn, word []byte) {
	fmt.Println("A client has connected!")
	fmt.Println("File send: ", string(word))
	nameFile := string(word)
	nameFile = strings.TrimFunc(nameFile, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	fmt.Println("File send: ", nameFile)
	// file, err := os.Open("6.png")
	file, err := os.Open(nameFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager1 *ClientManager1) start1() {
	for {
		select {
		case connection := <-manager1.register:
			manager1.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager1.unregister:
			if _, ok := manager1.clients[connection]; ok {
				close(connection.data)
				delete(manager1.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		case message := <-manager1.broadcast:
			for connection := range manager1.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager1.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		// go sendFileToClient(client.socket)
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED 1: " + string(message))
			manager.broadcast <- message
		}
	}
}

func (manager1 *ClientManager1) receive1(client1 *Client) {
	for {
		// go sendFileToClient(client.socket)
		message := make([]byte, 4096)
		length, err := client1.socket.Read(message)
		if err != nil {
			manager1.unregister <- client1
			client1.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED 1: " + string(message))
			manager1.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()

	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
			fmt.Println("Client connected")
			go sendFileToClient(client.socket, []byte(message))
		}
	}
}

func (manager1 *ClientManager1) send1(client1 *Client) {
	defer client1.socket.Close()

	for {
		select {
		case message, ok := <-client1.data:
			if !ok {
				return
			}
			client1.socket.Write(message)
			fmt.Println("Client connected")
			go sendFileToClient(client1.socket, []byte(message))
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":3000")
	if error != nil {
		fmt.Println(error)
	}
	listener1, error := net.Listen("tcp", ":3001")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	manager1 := ClientManager1{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager1.start1()

	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)

		connection1, _ := listener1.Accept()
		if error != nil {
			fmt.Println(error)
		}

		client1 := &Client{socket: connection1, data: make(chan []byte)}
		manager1.register <- client1

		go manager1.receive1(client1)
		go manager1.send1(client1)

	}
}

func main() {
	startServerMode()
}
