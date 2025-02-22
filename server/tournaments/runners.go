package tournaments

import (
	"log"
	"shared/interfaces"
	"sync"
	"tournament_server/models"

	pb "shared/grpc"
)

type MatchNotification struct {
	MatchId int
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
	currentMatches := r.tournament.PendingMatches()
	for  match := range currentMatches {
		r.runMatch( r.tournament.Matches()[match])
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
			if(r.tournament.Matches()[matchId].Status() == 2){
				r.tournament.Matches()[matchId].SetStatus(3)
				delete(r.tournament.PendingMatches(), matchId)
			}else{
				continue
			}
			winner := value.(interfaces.Player)
			r.tournament.Matches()[matchId].SetWinner(winner)
			match:= r.tournament.Matches()[matchId]
			log.Default().Printf("tournamentRunner:Notify: winner is %s\n", winner.Id())

			if matchId != 0 {
				r.tournament.Matches()[(matchId+1)/2-1].AddPlayer(winner)
				log.Default().Printf("tournamentRunner:Notify: setting winner %s for next match of %d\n", winner.Id(), matchId)

				nextMatchId := (matchId+1)/2-1
				log.Default().Printf("tournamentRunner:Notify: next match of %d is %d with %d players\n", match.Id(), nextMatchId, len(r.tournament.Matches()[(matchId+1)/2-1].Players()))
				if r.tournament.Matches()[(matchId+1)/2-1].SetGameIfReady() {
					log.Default().Printf("tournamentRunner:Notify: next match %d of %d has 2 players\n", nextMatchId, match.Id())
					r.tournament.PendingMatches()[nextMatchId] = true
					r.runMatch(r.tournament.Matches()[(matchId+1)/2-1])
					log.Default().Printf("tournamentRunner:Notify: running next match %d\n", nextMatchId)
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
	log.Default().Printf("matchRunner:Run: running match %d\n", r.match.Id())
	winner := r.match.Start()
	log.Default().Printf("matchRunner:Run: match %d finished\n", r.match.Id())

	channel <- MatchNotification{r.match.Id(), "finished", winner}
	log.Default().Printf("matchRunner:Run: winner is %s\n", winner.Id())
	r.runner.wg.Done()
}
