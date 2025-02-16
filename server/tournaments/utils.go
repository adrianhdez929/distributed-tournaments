package tournaments

import (
	"fmt"
	"tournament_server/models"
)

func GetStatistics(tournament models.Tournament) map[string]interface{} {
	playerWins := make(map[string]int32)
	// theres no winner yet
	finalWinner := tournament.Winner()

	for _, player := range tournament.Players() {
		playerWins[player.Id()] = 0
	}

	for k := range tournament.Matches() {
		stateKey := fmt.Sprintf("match_winner_%s", k)
		player := tournament.State()[stateKey]
		if player == nil {
			continue
		}

		playerWins[player.(string)]++
	}

	return map[string]interface{}{
		"player_wins": playerWins,
		"winner":      finalWinner,
	}
}
