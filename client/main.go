package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Read the content of players/greedy.go
	content, err := os.ReadFile("./client/players/greedy.go")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Connect to the server
	for {

		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			continue
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
}
