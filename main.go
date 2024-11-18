package main

import (
	// "tournaments/server/models"
	"fmt"
	"os"
)

// TODO: create both proyects: client and server

func main() {
	// playerCount := 16
	// players := make([]models.Player, playerCount)
	// // matches := make([]models.Match, playerCount/2)

	// for i := 0; i < playerCount; i++ {
	// 	players[i] = models.NewPlayerData(i + 1)
	// }

	// // for i := 1; i < playerCount; i += 2 {
	// // 	matches[i/2] = models.NewMatchData(players[i], players[i-1])
	// // }

	// tournament := models.NewTournamentData(players)

	// winner := tournament.Winner()

	// fmt.Printf("winner: %s\n", winner.Id())
	// server.Run()

	// _, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	fmt.Println("Port is in use")
	// }
	// client.Run()

	content, err := os.ReadFile("./server/players/greedy.go")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	server.ReceivePlayerImpl(string(content), "NewGreedyPlayer")
}
