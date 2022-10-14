package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

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

	gc.fs.StringVar(&gc.file, "file", "1", "channel of connect to be recieve")

	return gc
}

type ReceiveCommand struct {
	fs *flag.FlagSet

	channel string
}

type SendCommand struct {
	fs *flag.FlagSet

	file string
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
	// fmt.Println("Hello", g.channel, "!")
	// fmt.Println("Init app")
	// // Se conecta al server

	// serverCnn, err := net.Dial("tcp", ":8080")
	// if err != nil {
	// 	panic(err)
	// }

	// // escuchar
	// readMessages(serverCnn)
	if g.channel == "1" {
		fmt.Println("Hello", g.channel, "!")
		fmt.Println("Init app")
		// Se conecta al server

		serverCnn, err := net.Dial("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		// escuchar
		readMessages(serverCnn)
	}
	if g.channel == "2" {
		go receiveFile()

		var input string
		fmt.Scanln(&input)
	}

	return nil

}

func (g *SendCommand) Run() error {
	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		// Rompemos el ciclo
	}
	connection.Write([]byte("Hola mundo"))
	connection.Close()
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

func readMessages(cnn net.Conn) {
	// read message
	var message = make([]byte, 2048)

	_, err := cnn.Read(message)
	if err != nil {
		panic(err)
	}

	// print message
	fmt.Println("El mensaje del servidor es: ", string(message))
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
