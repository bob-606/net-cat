package main

import (
	"chat/history"
	"fmt"
	"sync"
	"time"
)

const col_reset = "\033[0m"

// project mandated Mutex
var messageLock = sync.Mutex{}

// sends the message to every client in ClientList
func BroadcastAll(message string) {
	messageLock.Lock()
	defer messageLock.Unlock()

	history.Write(message) // This logs the message to the file

	fmt.Print(message)

	asBytes := []byte(message)
	for _, c := range ClientList {
		c.Connection.Write(asBytes)
	}
}

// sends the message to each client in ClientList, except `client`
func BroadcastExcept(message string, client Client) {
	messageLock.Lock()
	defer messageLock.Unlock()

	history.Write(message)
	fmt.Print(message)

	asBytes := []byte(message)
	for _, c := range ClientList {
		if c == client {
			continue
		}

		c.Connection.Write(asBytes)
	}
}

func getTimeStamp() string {
	return fmt.Sprintf("[%v]", time.Now().Format("06/01/02 15:04"))
}

func ClientJoined(client Client) {
	message := fmt.Sprintf("%v %s joined the chat!\n", getTimeStamp(), client.Name)
	BroadcastAll(message)
}

func ClientLeft(client Client) {
	message := fmt.Sprintf("%v %s left the chat!\n", getTimeStamp(), client.Name)
	BroadcastAll(message)
}

func formatClient(client Client) string {
	return client.Name + ": "
}

func HandleMessage(client Client, message string) {
	if message == "" {
		return
	}

	// reconsider.
	formatted := fmt.Sprintf("%v %v %v%v\n", getTimeStamp(), formatClient(client), message, col_reset)

	// write to file aswell
	BroadcastExcept(formatted, client)
}

func sendHistoryToClient(client Client) {
	connection := client.Connection
	// write history to client
	connection.Write([]byte("\n-----------------------------------HISTORY--------\n"))
	for _, msg := range history.Get() {
		connection.Write([]byte(msg))
	}
	connection.Write([]byte("-----------------------------------HISTORY--------\n\n"))
}
