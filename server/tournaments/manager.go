package tournaments

import (
	"tournament_server/models"
)

type TournamentManager struct {
	Tournaments map[string]models.Tournament
}

func NewTournamentManager() *TournamentManager {
	return &TournamentManager{
		Tournaments: make(map[string]models.Tournament),
	}
}

func (tm *TournamentManager) AddTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament
}

func (tm *TournamentManager) UpdateTournament(tournament models.Tournament) {
	tm.Tournaments[tournament.Id()] = tournament
}
