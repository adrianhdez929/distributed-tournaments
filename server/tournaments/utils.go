package tournaments

import (
	"fmt"
	"shared/interfaces"
	"tournament_server/models"

	pb "shared/grpc"
)

func GetStatistics(tournament models.Tournament) map[string]interface{} {
	playerWins := make(map[string]int32)
	// theres no winner yet
	finalWinner := tournament.Winner()

	for _, player := range tournament.Players() {
		playerWins[player.Id()] = 0
	}

	for k := range tournament.Matches() {
		stateKey := fmt.Sprintf("match_winner_%d", k)
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

func DumpTournamentPlayers(players []interfaces.Player) []*pb.Player {
	dumpedPlayers := make([]*pb.Player, len(players))

	for _, player := range players {
		dumpedPlayers = append(
			dumpedPlayers,
			&pb.Player{
				Id:        player.Id(),
				Name:      player.Id(),
				AgentName: player.Name(),
			},
		)
	}

	return dumpedPlayers
}

func DumpTournamentMatches(m []models.Match) []*pb.Match {
	dumpedMatches := make([]*pb.Match, 0)

	for _, match := range m {
		if match == nil || match.Players() == nil {
			continue
		}

		players := DumpTournamentPlayers(match.Players())
		winner := match.Winner()
		next := ""

		dumpedMatches = append(
			dumpedMatches,
			&pb.Match{
				Id:      fmt.Sprintf("%d", match.Id()),
				Player1: players[0],
				Player2: players[0],
				Winner: &pb.Player{
					Id:        winner.Id(),
					Name:      winner.Id(),
					AgentName: winner.Name(),
				},
				Next: next,
			},
		)
	}

	return dumpedMatches
}
