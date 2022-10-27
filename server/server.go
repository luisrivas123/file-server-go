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

var countClientsChannel1 = 0
var countClientsChannel2 = 0
var countClientsChannel3 = 0

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

func sendFileToClient(connection net.Conn, word []byte) {
	// fmt.Println("File send: ", string(word))
	nameFile := string(word)
	nameFile = strings.TrimFunc(nameFile, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
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
	// fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	// fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	// fmt.Println("File has been sent, closing connection!")
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

func (manager *ClientManager) start1() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			countClientsChannel1++
			fmt.Println("Added new connection channel 1! clients: ", countClientsChannel1)
			// fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				countClientsChannel1--
				fmt.Println("A connection has terminated channel 1! clients: ", countClientsChannel1)
				// fmt.Println("A connection has terminated!")
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

func (manager *ClientManager) start2() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			countClientsChannel2++
			fmt.Println("Added new connection channel 2! clients: ", countClientsChannel2)
			// fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				countClientsChannel2--
				fmt.Println("A connection has terminated channel 2! clients: ", countClientsChannel2)
				// fmt.Println("A connection has terminated!")
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

func (manager *ClientManager) start3() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			countClientsChannel3++
			fmt.Println("Added new connection channel 3! clients: ", countClientsChannel3)
			// fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				countClientsChannel3--
				fmt.Println("A connection has terminated channel 3! clients: ", countClientsChannel3)
				// fmt.Println("A connection has terminated!")
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

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			// fmt.Println("RECEIVED : " + string(message))
			manager.broadcast <- message
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
			// fmt.Println("Client connected")
			go sendFileToClient(client.socket, []byte(message))
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server ...")
	go channel_1()
	go channel_2()
	go channel_3()
}

func main() {
	startServerMode()
	var input string
	fmt.Scanln(&input)
}

func channel_1() {
	// fmt.Println("channel one open...")
	listener, error := net.Listen("tcp", ":3000")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start1()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)

	}
}

func channel_2() {
	// fmt.Println("channel two open...")
	listener, error := net.Listen("tcp", ":3001")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start2()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)

	}
}

func channel_3() {
	// fmt.Println("channel three open...")
	listener, error := net.Listen("tcp", ":3002")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start3()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)

	}
}
