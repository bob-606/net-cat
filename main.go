package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

var shutdown = false
var listener net.Listener
var logger *log.Logger

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		log.Fatal("Expecting 1 <= arguments for port number")
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	}

	address := fmt.Sprintf(":%v", port)
	fmt.Printf("Listening on port %v\n", address)

	go waitShutdownRequest()

	var server_err error
	listener, server_err = net.Listen("tcp", address)
	if server_err != nil {
		log.Fatal(server_err)
	}

	for {
		connection, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return // sure is a workaround.
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		go NewClient(connection)
	}
}

func waitShutdownRequest() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type `f` + Enter to shut down the server")
	for {
		scanner.Scan()
		txt := scanner.Text()
		if txt == "f" {
			break
		}
	}

	listener.Close()
	os.Exit(0)
}
