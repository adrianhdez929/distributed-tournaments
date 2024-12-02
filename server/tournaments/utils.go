package tournaments

import (
	"fmt"
	"tournament_server/models"
)

func GetStatistics(tournament models.Tournament) map[string]interface{} {
	playerWins := make(map[string]int32)
	// theres no winner yeu
	finalWinner := tournament.State()["winner"]

	if finalWinner != nil {
		finalWinner = finalWinner.(string)
	} else {
		finalWinner = ""
	}

	for _, player := range tournament.Players() {
		playerWins[player.Id()] = 0
	}

	for i := range tournament.Matches() {
		stateKey := fmt.Sprintf("match_winner_%d", i)
		player := tournament.State()[stateKey]
		if player == nil {
			continue
		}

		playerWins[player.(string)]++
	}

	return map[string]interface{}{
		"player_wins":  playerWins,
		"final_winner": finalWinner,
	}
}
