package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

type Client struct {
	socket net.Conn
	data   chan []byte
}

func NewReceiveCommand() *ReceiveCommand {
	gc := &ReceiveCommand{
		fs: flag.NewFlagSet("receive", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.channel, "channel", "1", "channel of connect to be recieve")

	return gc
}

func NewSendCommand() *SendCommand {
	gc := &SendCommand{
		fs: flag.NewFlagSet("send", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.channel, "channel", "1", "channel of connect to be recieve")

	return gc
}

type ReceiveCommand struct {
	fs *flag.FlagSet

	channel string
}

type SendCommand struct {
	fs *flag.FlagSet

	channel string
}

func (g *ReceiveCommand) Name() string {
	return g.fs.Name()
}

func (g *SendCommand) Name() string {
	return g.fs.Name()
}

func (g *ReceiveCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *SendCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ReceiveCommand) Run() error {
	if g.channel == "1" {
		fmt.Println("Starting client recieve channel 1...")
		connection, error := net.Dial("tcp", "localhost:3000")
		if error != nil {
			fmt.Println(error)
		}
		defer connection.Close()

		client := &Client{socket: connection}
		go client.receive()

		var input string
		fmt.Scanln(&input)
	}
	if g.channel == "2" {
		fmt.Println("Starting client recieve channel  2...")
		connection, err := net.Dial("tcp", "localhost:3001")
		if err != nil {
			panic(err)
		}
		defer connection.Close()

		client := &Client{socket: connection}
		go client.receive()

		var input string
		fmt.Scanln(&input)
	}
	if g.channel == "3" {
		fmt.Println("Starting client recieve channel  3...")
		connection, err := net.Dial("tcp", "localhost:3002")
		if err != nil {
			panic(err)
		}
		defer connection.Close()

		client := &Client{socket: connection}
		go client.receive()

		var input string
		fmt.Scanln(&input)
	}
	return nil

}

func (g *SendCommand) Run() error {
	if g.channel == "1" {
		fmt.Println("Starting client send channel 1...")
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
	if g.channel == "2" {
		fmt.Println("Starting client send channel 2...")
		connection, error := net.Dial("tcp", "localhost:3001")
		if error != nil {
			fmt.Println(error)
		}
		for {
			reader := bufio.NewReader(os.Stdin)
			message, _ := reader.ReadString('\n')
			connection.Write([]byte(strings.TrimRight(message, "\n")))
		}
	}
	if g.channel == "3" {
		fmt.Println("Starting client send channel 3...")
		connection, error := net.Dial("tcp", "localhost:3002")
		if error != nil {
			fmt.Println(error)
		}
		for {
			reader := bufio.NewReader(os.Stdin)
			message, _ := reader.ReadString('\n')
			connection.Write([]byte(strings.TrimRight(message, "\n")))
		}
	}
	return nil
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		NewReceiveCommand(),
		NewSendCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {

	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		receiveFile(client.socket)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

func receiveFile(con net.Conn) {
	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	con.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	con.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	// path := filepath.Join("./file/", fileName)

	newFile, err := os.Create(fileName)
	// newFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, con, (fileSize - receivedBytes))
			con.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, con, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}
