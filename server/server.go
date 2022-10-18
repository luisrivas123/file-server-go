package main

import (
	"fmt"
	"io"
	"net"
	"os"
	// "strconv"
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
			fmt.Println("RECEIVED 1: " + string(message))
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
			sendFile(client.socket, string(message))
		}
	}
}

func startServerMode() {
	fmt.Println("Starting server...")
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
	go manager.start()
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

func main() {
	startServerMode()
}

// func sendFile(server net.Listener, srcFile string) {
func sendFile(connection net.Conn, srcFile string) {
	// accept connection
	// conn, err := server.Accept()
	// if err != nil {
	// 	fmt.Println("Fatal error: ", err)
	// }

	// open file to send
	// fi, err := os.Open(srcFile)
	fi, err := os.Open("./file/6.png")
	if err != nil {
		fmt.Println("Fatal error: ", err)
	}
	// defer fi.Close()

	// send file to client
	_, err = io.Copy(connection, fi)
	if err != nil {
		fmt.Println("Fatal error: ", err)
	}
}