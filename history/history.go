package history

import (
	"bufio"
	"log"
	"os"
)

const historyLen = 5

var logFile *os.File
var activeHistory []string

func init() {
	//adding logging functionality
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logFile = file

	// TODO: improve file read
	// currently just reading every single line of the file and pushing each onto activeHistory
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		addToActive(scanner.Text() + "\n")
	}
}

func addToActive(message string) {
	if len(activeHistory) > historyLen {
		activeHistory = activeHistory[1:]
	}

	activeHistory = append(activeHistory, message)
}

func Write(message string) {
	addToActive(message)
	logFile.Write([]byte(message))
}

func Get() []string {
	return activeHistory
}
