package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Read the content of players/greedy.go
	content, err := os.ReadFile("./players/greedy.go")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Connect to the server

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Send the content through the socket
	_, err = conn.Write(content)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("File content sent successfully")
}
