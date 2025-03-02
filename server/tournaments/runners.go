package tournaments

import (
	"log"
	pb "shared/grpc"
	"shared/interfaces"
	"sync"
	"time"
	"tournament_server/models"
)

type MatchNotification struct {
	MatchId int
	Key     string
	Value   interface{}
}

type TournamentRunner struct {
	manager    *TournamentManager
	tournament models.Tournament
	lock       *sync.Mutex
	wg         *sync.WaitGroup
	done       bool
}

func NewTournamentRunner(manager *TournamentManager, tournament models.Tournament) *TournamentRunner {
	return &TournamentRunner{
		manager:    manager,
		tournament: tournament,
		lock:       &sync.Mutex{},
		wg:         &sync.WaitGroup{},
		done:       false,
	}
}

func (r *TournamentRunner) runMatch(match models.Match, isLast bool, channel chan MatchNotification) {
	if match.SetGameIfReady() {
		matchRunner := NewMatchRunner(r, match, isLast)
		matchRunner.Run(channel)
	}
}

func (r *TournamentRunner) Run() {
	log.Default().Printf("tournamentRunner:Run: tournament %s started\n", r.tournament.Id())
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS)

	notifyChannel := make(chan MatchNotification)

	// Distribuir la ejecucion de las partidas tambien. habra que agregar un layer nuevo para esto? o sobre este mismo no se pierde flexibilidad?
	currentMatches := r.tournament.PendingMatches()
	for match := range currentMatches {
		time.Sleep(3 * time.Second)

		if r.tournament.Matches()[match].Id() == len(r.tournament.Players())-1 {
			go r.runMatch(r.tournament.Matches()[match], true, notifyChannel)
		} else {
			go r.runMatch(r.tournament.Matches()[match], false, notifyChannel)
		}
	}

	r.Notify(notifyChannel)

	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_COMPLETED)

	log.Default().Printf("tournamentRunner:Run: tournament %s finished\n", r.tournament.Id())
}

func (r *TournamentRunner) Resume() {
	log.Default().Printf("tournamentRunner:Resume: tournament %s resumed\n", r.tournament.Id())
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS)

	notifyChannel := make(chan MatchNotification)

	currentMatches := r.tournament.PendingMatches()
	for match := range currentMatches {
		time.Sleep(3 * time.Second)

		if r.tournament.Matches()[match].Id() == len(r.tournament.Players())-1 {
			go r.runMatch(r.tournament.Matches()[match], true, notifyChannel)
		} else {
			go r.runMatch(r.tournament.Matches()[match], false, notifyChannel)
		}
	}

	r.Notify(notifyChannel)

	log.Printf("tournamentRunner:Resume: out of notify\n")
	r.manager.Notify(r.tournament.Id(), "status", pb.TournamentStatus_TOURNAMENT_STATUS_COMPLETED)

	log.Default().Printf("tournamentRunner:Run: tournament %s finished\n", r.tournament.Id())
}

func (r *TournamentRunner) Notify(notifyChannel chan MatchNotification) {
	log.Printf("tournamentRunnser: Notify entering notify loop")
	for notification := range notifyChannel {
		matchId := notification.MatchId
		key := notification.Key
		value := notification.Value

		log.Printf("tournamentRunner:Notify: received notification with matchId %d key %s and value %v\n", matchId, key, value)

		if key == "finished" {
			if r.tournament.Matches()[matchId].Status() == 1 {
				r.lock.Lock()
				r.tournament.Matches()[matchId].SetStatus(3)
				delete(r.tournament.PendingMatches(), matchId)
				r.manager.UpdateTournament(r.tournament)
				r.lock.Unlock()
			} else {
				r.lock.Lock()
				r.manager.UpdateTournament(r.tournament)
				r.lock.Unlock()
				continue
			}

			winner := value.(interfaces.Player)
			r.lock.Lock()
			r.tournament.Matches()[matchId].SetWinner(winner)
			match := r.tournament.Matches()[matchId]
			r.lock.Unlock()
			log.Default().Printf("tournamentRunner:Notify: winner of %d is %s\n", matchId, winner.Id())

			if matchId != 0 {
				nextMatchId := (matchId+1)/2 - 1

				log.Default().Printf("tournamentRunner:Notify: setting winner %s for next match of %d\n", winner.Id(), matchId)
				r.lock.Lock()
				r.tournament.Matches()[nextMatchId].AddPlayer(winner)
				r.lock.Unlock()
				log.Default().Printf("tournamentRunner:Notify: next match of %d is %d with %d players\n", match.Id(), nextMatchId, len(r.tournament.Matches()[nextMatchId].Players()))

				if r.tournament.Matches()[nextMatchId].SetGameIfReady() {
					log.Default().Printf("tournamentRunner:Notify: next match %d of %d has 2 players\n", nextMatchId, match.Id())
					r.lock.Lock()
					r.tournament.PendingMatches()[nextMatchId] = true
					r.lock.Unlock()
					time.Sleep(3 * time.Second)

					if nextMatchId == len(r.tournament.Players())-1 {
						log.Printf("tournamentRunner:Notify: detected last match with id %d", nextMatchId)
						go r.runMatch(r.tournament.Matches()[nextMatchId], true, notifyChannel)
					} else {
						go r.runMatch(r.tournament.Matches()[nextMatchId], false, notifyChannel)
					}
					log.Default().Printf("tournamentRunner:Notify: running next match %d\n", nextMatchId)
					r.manager.UpdateTournament(r.tournament)
				}
			} else {
				r.manager.Notify(r.tournament.Id(), "finished", winner)
				log.Default().Printf("tournamentRunner:Notify: tournament winner is %s\n", winner.Id())
				r.manager.UpdateTournament(r.tournament)
				break
			}
		}
	}
}

func (r *TournamentRunner) Save() {
	tournament_json := r.tournament.ToJson()
	r.manager.SaveTournamentAsJson(r.tournament.Id(), tournament_json)
}

type MatchRunner struct {
	runner      *TournamentRunner
	match       models.Match
	isLastMatch bool
}

func NewMatchRunner(
	runner *TournamentRunner,
	match models.Match,
	isLast bool,
) *MatchRunner {
	return &MatchRunner{
		runner:      runner,
		match:       match,
		isLastMatch: isLast,
	}
}

func (r *MatchRunner) Run(channel chan MatchNotification) {
	log.Default().Printf("matchRunner:Run: running match %d\n", r.match.Id())
	winner := r.match.Start()
	log.Default().Printf("matchRunner:Run: match %d finished\n", r.match.Id())

	channel <- MatchNotification{r.match.Id(), "finished", winner}
	log.Default().Printf("matchRunner:Run: winner is %s\n", winner.Id())
}
