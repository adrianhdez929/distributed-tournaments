package tournaments

import (
	"log"
	"shared/interfaces"
	"sync"
	"tournament_server/models"

	pb "shared/grpc"
)

type MatchNotification struct {
	MatchId string
	Key     string
	Value   interface{}
}

type TournamentRunner struct {
	manager       *TournamentManager
	tournament    models.Tournament
	wg            *sync.WaitGroup
	notifyChannel chan MatchNotification
	done          bool
}

func NewTournamentRunner(manager *TournamentManager, tournament models.Tournament) *TournamentRunner {
	return &TournamentRunner{
		manager:       manager,
		tournament:    tournament,
		wg:            &sync.WaitGroup{},
		notifyChannel: make(chan MatchNotification),
		done:          false,
	}
}

func (r *TournamentRunner) runMatch(match models.Match) {
	matchRunner := NewMatchRunner(r, match)
	r.wg.Add(1)
	go matchRunner.Run(r.notifyChannel)
}

func (r *TournamentRunner) Run() {
	log.Default().Printf("tournamentRunner:Run: tournament %s started\n", r.tournament.Id())
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS)

	go r.Notify()

	// Distribuir la ejecucion de las partidas tambien. habra que agregar un layer nuevo para esto? o sobre este mismo no se pierde flexibilidad?
	currentMatches := r.tournament.InitialMatches()
	for _, match := range currentMatches {
		r.tournament.SetMatch(match.Id(), match)
		r.runMatch(match)
	}

	r.wg.Wait()
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_COMPLETED)

	log.Default().Printf("tournamentRunner:Run: tournament %s finished\n", r.tournament.Id())
}

func (r *TournamentRunner) Notify() {
	r.wg.Add(1)

	for {
		notification := <-r.notifyChannel

		matchId := notification.MatchId
		key := notification.Key
		value := notification.Value

		if key == "finished" {
			winner := value.(interfaces.Player)
			log.Default().Printf("tournamentRunner:Notify: winner is %s\n", winner.Id())
			match := r.tournament.Matches()[matchId]
			match.SetWinner(winner)
			r.tournament.SetMatch(matchId, match)

			if r.tournament.Matches()[matchId].Next() != nil {
				r.tournament.Matches()[matchId].Next().SetPlayer(winner)
				log.Default().Printf("tournamentRunner:Notify: setting winner %s for next match of %s\n", winner.Id(), matchId)

				nextMatch := match.Next()
				log.Default().Printf("tournamentRunner:Notify: next match of %s is %s with %d players\n", match.Id(), nextMatch.Id(), len(nextMatch.Players()))
				if len(nextMatch.Players()) == 2 {
					log.Default().Printf("tournamentRunner:Notify: next match %s of %s has 2 players\n", nextMatch.Id(), match.Id())
					r.tournament.SetMatch(nextMatch.Id(), nextMatch)
					r.runMatch(nextMatch)
					log.Default().Printf("tournamentRunner:Notify: running next match %s\n", nextMatch.Id())
				}
			} else {
				r.manager.Notify(r.tournament.Id(), "finished", winner)
				log.Default().Printf("tournamentRunner:Notify: tournament winner is %s\n", winner.Id())
				break
			}
		}
	}

	r.wg.Done()
}

type MatchRunner struct {
	runner *TournamentRunner
	match  models.Match
}

func NewMatchRunner(
	runner *TournamentRunner,
	match models.Match,
) *MatchRunner {
	return &MatchRunner{
		runner: runner,
		match:  match,
	}
}

func (r *MatchRunner) Run(channel chan MatchNotification) {
	log.Default().Printf("matchRunner:Run: running match %s\n", r.match.Id())
	winner := r.match.Start()
	log.Default().Printf("matchRunner:Run: match %s finished\n", r.match.Id())

	channel <- MatchNotification{r.match.Id(), "finished", winner}
	log.Default().Printf("matchRunner:Run: winner is %s\n", winner.Id())
	r.runner.wg.Done()
}
