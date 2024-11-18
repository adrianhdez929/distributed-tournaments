package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"plugin"
	"shared/interfaces"
	"strings"
	"tournament_server/models"
)

func main() {
	// content, err := os.ReadFile("./players/greedy.go")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// ReceivePlayerImpl(string(content), "NewGreedyPlayer")
	Run()
}

func Run() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	handleConnection(listener)

	fmt.Println("Server is listening on port 8080")

}

func handleConnection(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		receivedString := string(buffer[:n])
		// fmt.Println("Received string:", receivedString)
		playerFactory, err := getPlayerConstructor(receivedString, "NewGreedyPlayer")
		if err != nil {
			fmt.Println("Error building dynamic object:", err)
			return
		}
		playerCount := 16

		players := make([]interfaces.Player, playerCount)
		// matches := make([]models.Match, playerCount/2)

		for i := 0; i < playerCount; i++ {
			players[i] = playerFactory(i + 1)
			fmt.Printf("creating player %s\n", players[i].Id())
		}

		// for i := 1; i < playerCount; i += 2 {
		// 	matches[i/2] = models.NewMatchData(players[i], players[i-1])
		// }

		tournament := models.NewTournamentData(players)

		winner := tournament.Winner()
		fmt.Printf("the winner is %s\n", winner.Id())

		break
	}
}

func getPlayerConstructor(code string, constructor string) (func(int) interfaces.Player, error) {
	filename := "./players/player_impl.go"
	file, _ := os.Create(filename)
	code = strings.Replace(code, "package players", "package main", 1)
	file.WriteString(code)
	file.Close()
	defer os.Remove(filename)

	// TODO: add unique identifier to players uploaded
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "impl.so", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return nil, err
	}
	if err != nil {
		fmt.Printf("error in exec: %s\n", err)
		return nil, err
	}
	plug, err := plugin.Open("impl.so")
	if err != nil {
		fmt.Printf("error plugin: %s\n", err)
		return nil, err
	}
	defer os.Remove("impl.so")
	playerImpl, err := plug.Lookup(constructor)
	if err != nil {
		fmt.Printf("error in lookup: %s\n", err)
		return nil, err
	}

	return playerImpl.(func(int) interfaces.Player), nil
}

func ReceivePlayerImpl(code string, constructor string) {
	file, _ := os.Create("player_impl.go")
	file.WriteString(code)
	file.Close()

	// TODO: add unique identifier to players uploaded
	_, err := exec.Command("go", "build", "-buildmode=plugin", "-o", "impl.so", "player_impl.go").Output()
	if err != nil {
		fmt.Printf("error in exec: %s\n", err)
	}
	plug, err := plugin.Open("impl.so")
	if err != nil {
		fmt.Printf("error plugin: %s\n", err)
	}
	playerImpl, err := plug.Lookup(constructor)
	if err != nil {
		fmt.Printf("error in lookup: %s\n", err)
	}
	// TODO: create shared interface for player definition
	loadedPlayer := playerImpl.(func(int) interfaces.Player)(1)

	loadedPlayer.Move()
}
