package tournaments

import (
	"context"
	"errors"
	"tournament_server/models"

	pb "shared/grpc"
)

type TournamentManager struct {
	Tournaments map[string]models.Tournament
	repo        TournamentRepository
}

func NewTournamentManager(repo TournamentRepository) *TournamentManager {
	return &TournamentManager{
		Tournaments: make(map[string]models.Tournament),
		repo:        repo,
	}
}

// func (tm *TournamentManager) SaveTournament(tournament models.Tournament) error {
// 	return tm.repo.(tournament)
// }

func (tm *TournamentManager) AddTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament

	go NewTournamentRunner(tm, tournament).Run()
}

func (tm *TournamentManager) UpdateTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament
}

func (tm *TournamentManager) GetTournament(id string) (models.Tournament, error) {
	tournament, ok := tm.Tournaments[id]
	if !ok {
		return nil, errors.New("tournament not found")
	}

	return tournament, nil
}

func (tm *TournamentManager) Notify(id string, key string, value interface{}) error {
	tournament := tm.Tournaments[id]

	if tournament == nil {
		return errors.New("tournament not found")
	}

	if key == "finished" {
		return nil
	}

	if key == "status" {
		tournament.SetStatus(value.(pb.TournamentStatus))
	}

	tournament.SetState(key, value)

	statistics := GetStatistics(tournament)

	pbTournament := &pb.Tournament{
		Id:          tournament.Id(),
		Status:      tournament.Status(),
		PlayerWins:  statistics["player_wins"].(map[string]int32),
		FinalWinner: statistics["final_winner"].(string),
	}

	tm.repo.Update(context.Background(), pbTournament)

	return nil
}
