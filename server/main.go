package main

import (
	"fmt"
	"net"
	"os"
	"shared/interfaces"
	"tournament_server/games"
	"tournament_server/models"
	"tournament_server/players"
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

		// TODO: make it a non anonymous function
		go func(c net.Conn) {
			defer conn.Close()

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				return
			}

			receivedString := string(buffer[:n])
			fmt.Println("Received string:", receivedString)
			// playerFactory, err := code.GetPlayerConstructor(receivedString, "NewGreedyPlayer")
			playerFactory := players.NewGreedyPlayer
			if err != nil {
				fmt.Println("Error building dynamic object:", err)
				return
			}

			// gameFactory, err := code.GetGameConstructor(receivedString, "NewTicTacToe")
			gameFactory := games.NewTicTacToe
			if err != nil {
				fmt.Println("Error building dynamic object:", err)
				return
			}

			createTournament(playerFactory, gameFactory, 16)
		}(conn)
	}
}

func createTournament(playerFactory func(int) interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game, playerCount int) {
	players := make([]interfaces.Player, playerCount)
	// matches := make([]models.Match, playerCount/2)

	for i := 0; i < playerCount; i++ {
		players[i] = playerFactory(i + 1)
		fmt.Printf("creating player %s\n", players[i].Id())
	}

	tournament := models.NewTournamentData(players, gameFactory)

	winner := tournament.Winner()
	fmt.Printf("the winner is %s\n", winner.Id())
}
