package tournaments

import (
	"context"
	"errors"
	"log"
	"tournament_server/chord"
	"tournament_server/models"

	pb "shared/grpc"
	"shared/interfaces"
)

type TournamentManager struct {
	Tournaments map[string]models.Tournament
	repo        TournamentRepository
	node        *chord.ChordServer
}

func NewTournamentManager(repo TournamentRepository, server *chord.ChordServer) *TournamentManager {
	return &TournamentManager{
		Tournaments: make(map[string]models.Tournament),
		repo:        repo,
		node:        server,
	}
}

// func (tm *TournamentManager) SaveTournament(tournament models.Tournament) error {
// 	return tm.repo.(tournament)
// }

func (tm *TournamentManager) AddTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament

	go NewTournamentRunner(tm, tournament).Run()
	// Aqui hay que hacer la replicacion a los siguientes k - 1 nodos con k factor de replicacion
}

func (tm *TournamentManager) UpdateTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament
}

func (tm *TournamentManager) GetTournament(id string) (models.Tournament, error) {
	// id = name en este caso
	tournament, ok := tm.Tournaments[id]
	if !ok {
		return nil, errors.New("tournament not found")
	}

	return tournament, nil
}

func (tm *TournamentManager) Notify(tournamentId string, key string, value interface{}) error {
	// Revisar esto, ya que va a ser el core del sistema de notificaciones del sistema
	tournament := tm.Tournaments[tournamentId]

	if tournament == nil {
		return errors.New("tournament not found")
	}

	if key == "finished" {
		winner := value.(interfaces.Player)
		tournament.SetWinner(winner)

		tm.Tournaments[tournamentId] = tournament

		statistics := GetStatistics(tournament)

		pbTournament := &pb.Tournament{
			Id:              tournament.Id(),
			Name:            tournament.Id(),
			Status:          tournament.Status(),
			MaxParticipants: int32(len(tournament.Players())),
			Game:            tournament.Game(),
			Players:         DumpTournamentPlayers(tournament.Players()),
			Matches:         DumpTournamentMatches(tournament.Matches()),
			PlayerWins:      statistics["player_wins"].(map[string]int32),
			FinalWinner:     statistics["winner"].(interfaces.Player).Id(),
		}

		tm.repo.Update(context.Background(), pbTournament)

		log.Default().Printf("tournamentManager:Notify: tournament %s winner is %s\n", tournament.Id(), tournament.Winner().Id())
		return nil
	}

	if key == "status" {
		tournament.SetStatus(value.(pb.TournamentStatus))
		tm.Tournaments[tournamentId] = tournament

		statistics := GetStatistics(tournament)

		pbTournament := &pb.Tournament{
			Id:              tournament.Id(),
			Name:            tournament.Id(),
			Status:          tournament.Status(),
			MaxParticipants: int32(len(tournament.Players())),
			Game:            tournament.Game(),
			Players:         DumpTournamentPlayers(tournament.Players()),
			Matches:         DumpTournamentMatches(tournament.Matches()),
			PlayerWins:      statistics["player_wins"].(map[string]int32),
			FinalWinner:     "",
		}

		tm.repo.Update(context.Background(), pbTournament)

		log.Default().Printf("tournamentManager:Notify: tournament %s status is %s\n", tournament.Id(), tournament.Status())
	}

	return nil
}

// Hay que crear un worker que cuando finalice una partida/torneo, notifique a los nodos que puedan tener la partida
// para que actualicen/detengan la ejecucion
