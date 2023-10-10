package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var ServerCapacity = 10
var ClientList = []Client{}

type Client struct {
	Name       string
	Connection net.Conn
}

func NewClient(connecion net.Conn) {
	client := createClient(connecion)
	defer removeClient(client)

	if len(ClientList) >= ServerCapacity {
		connecion.Write([]byte("Cant join, server is full, sorry!!"))
		connecion.Close()
		return
	}

	sendHistoryToClient(client)

	ClientList = append(ClientList, client)
	ClientJoined(client)

	scanner := bufio.NewScanner(connecion)
	for scanner.Scan() {
		HandleMessage(client, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
}

func createClient(connection net.Conn) (client Client) {
	client.Connection = connection

	// Read greeting from the file
	// Consider whether you care about this been changeable at runtime
	// or mb just have a simple text backup
	greeting, err := os.ReadFile("greeting.txt")
	if err != nil {
		connection.Write([]byte("Error: " + err.Error()))
		connection.Close()
		panic(err)
	}

	connection.Write(greeting)

	// Read client's name
	scanner := bufio.NewScanner(connection)
	for {
		scanner.Scan()
		name := scanner.Text()

		// verify name
		if len(name) >= 3 {
			client.Name = name
			break
		}
		connection.Write([]byte("atleast 3 characters\n"))
	}

	return client
}

func removeClient(client Client) {
	client.Connection.Close()

	c_index := -1
	for i, c := range ClientList {
		if c == client {
			c_index = i
			break
		}
	}
	if c_index != -1 {
		ClientList = append(ClientList[:c_index], ClientList[c_index+1:]...)
	}

	// Notify others that this client has left the chat
	ClientLeft(client)
	fmt.Println("Exit:", client.Name, c_index)
}
