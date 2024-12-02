package tournaments

import (
	"fmt"
	"tournament_server/models"

	pb "shared/grpc"
)

type TournamentRunner struct {
	manager    *TournamentManager
	tournament models.Tournament
}

func NewTournamentRunner(manager *TournamentManager, tournament models.Tournament) *TournamentRunner {
	return &TournamentRunner{
		manager:    manager,
		tournament: tournament,
	}
}

func (r *TournamentRunner) Run() {
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS)

	currentMatches := r.tournament.Matches()

	for i := 0; i < r.tournament.TotalRounds(); i++ {
		newMatches := make([]models.Match, 0, len(currentMatches)/2)

		for _, v := range currentMatches {
			v.Start()
			winner := v.Winner()

			if v.Next() != nil {
				v.Next().SetPlayer(winner)
				r.manager.Notify(r.tournament.Id(), fmt.Sprintf("match_winner_%d", i), winner.Id())
			} else {
				r.manager.Notify(r.tournament.Id(), "winner", winner.Id())
			}

			if len(newMatches) == 0 {
				newMatches = append(newMatches, v.Next())
			} else if newMatches[len(newMatches)-1] != v.Next() {
				newMatches = append(newMatches, v.Next())
			}
		}

		currentMatches = newMatches
		r.tournament.AddMatches(newMatches)
	}

	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_COMPLETED)
}
