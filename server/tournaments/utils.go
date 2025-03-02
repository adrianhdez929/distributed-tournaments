package tournaments

import (
	"fmt"
	"shared/interfaces"
	"tournament_server/chord"
	"tournament_server/models"

	pb "shared/grpc"
)

func GetStatistics(tournament models.Tournament) map[string]interface{} {
	playerWins := make(map[string]int32)
	// theres no winner yet
	finalWinner := tournament.Winner()

	for _, player := range tournament.Players() {
		playerWins[player.Id()] = 0
		if finalWinner != nil && player.Id() == finalWinner.Id() {
			playerWins[player.Id()]++
		}
	}

	for _, m := range tournament.Matches() {
		player := m.Winner()

		if player == nil {
			continue
		}

		playerWins[player.Id()]++
	}

	return map[string]interface{}{
		"player_wins": playerWins,
		"winner":      finalWinner,
	}
}

func DumpTournamentPlayers(players []interfaces.Player) []*pb.Player {
	dumpedPlayers := make([]*pb.Player, 0)

	for _, player := range players {
		if player.Id() != "" {
			dumpedPlayers = append(
				dumpedPlayers,
				&pb.Player{
					Id:        player.Id(),
					Name:      player.Id(),
					AgentName: player.Name(),
				},
			)
		}
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

		var winnerPb *pb.Player

		if winner != nil {
			winnerPb = &pb.Player{
				Id:        winner.Id(),
				Name:      winner.Id(),
				AgentName: winner.Name(),
			}
		}

		if len(players) == 1 {
			dumpedMatches = append(
				dumpedMatches,
				&pb.Match{
					Id:      fmt.Sprintf("%d", match.Id()),
					Player1: players[0],
					Winner:  winnerPb,
					Next:    next,
				},
			)
		} else if len(players) == 2 {
			dumpedMatches = append(
				dumpedMatches,
				&pb.Match{
					Id:      fmt.Sprintf("%d", match.Id()),
					Player1: players[0],
					Player2: players[1],
					Winner:  winnerPb,
					Next:    next,
				},
			)
		} else {
			dumpedMatches = append(
				dumpedMatches,
				&pb.Match{
					Id:     fmt.Sprintf("%d", match.Id()),
					Winner: winnerPb,
					Next:   next,
				},
			)
		}
	}

	return dumpedMatches
}

func DumpTournament(tournament models.Tournament) *pb.Tournament {
	statistics := GetStatistics(tournament)
	winner := ""

	if tournament.Winner() != nil {
		winner = tournament.Winner().Id()
	}

	return &pb.Tournament{
		Id:              tournament.Id(),
		Name:            tournament.Id(),
		Status:          tournament.Status(),
		MaxParticipants: int32(len(tournament.Players())),
		Game:            tournament.Game(),
		Players:         DumpTournamentPlayers(tournament.Players()),
		Matches:         DumpTournamentMatches(tournament.Matches()),
		PlayerWins:      statistics["player_wins"].(map[string]int32),
		FinalWinner:     winner,
	}
}

func GetTournamentKey(name string) string {
	return fmt.Sprintf("tournament:%s", name)
}

func GetTournamentOwner(node *chord.ChordServer, name string) chord.ChordNodeReference {
	tHash := node.GetSha(GetTournamentKey(name))
	owner := node.FindSuccessor(tHash)
	return owner
}
