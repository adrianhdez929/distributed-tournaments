package tournaments

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"tournament_server/chord"
	"tournament_server/models"

	pb "shared/grpc"
	"shared/interfaces"
)

type TournamentManager struct {
	Tournaments map[string]models.Tournament
	repo        TournamentRepository
	node        *chord.ChordServer
	lock        *sync.Mutex
}

func NewTournamentManager(repo TournamentRepository, server *chord.ChordServer) *TournamentManager {
	return &TournamentManager{
		Tournaments: make(map[string]models.Tournament),
		repo:        repo,
		node:        server,
		lock:        &sync.Mutex{},
	}
}

// func (tm *TournamentManager) SaveTournament(tournament models.Tournament) error {
// 	return tm.repo.(tournament)
// }

func (tm *TournamentManager) StartTournament(tournament models.Tournament) {
	tm.lock.Lock()
	tm.Tournaments[tournament.Id()] = tournament
	tm.lock.Unlock()

	tm.repo.Create(context.Background(), DumpTournament(tournament))

	json, err := TournamentToJson(DumpTournament(tournament))
	if err != nil {
		log.Printf("manager:StartTournament: error while serializing tournament data: %s", err)
		return
	}

	tm.node.StoreKey(GetTournamentKey(tournament.Id()), json, chord.REPLICATION_FACTOR, chord.UPDATE)
	go NewTournamentRunner(tm, tournament).Run()
	// Aqui hay que hacer la replicacion a los siguientes k - 1 nodos con k factor de replicacion
}

func (tm *TournamentManager) ResumeTournament(json string) {
	tournament := (&models.TournamentData{}).FromJson(json)
	if tournament == nil {
		log.Printf("manager:ResumeTournament: cannot load tournament from json")
	}

	tm.lock.Lock()
	tm.Tournaments[tournament.Id()] = tournament
	tm.lock.Unlock()

	tournamentPb := DumpTournament(tournament)
	tm.repo.Create(context.Background(), tournamentPb)

	data, err := TournamentToJson(tournamentPb)
	if err != nil {
		log.Printf("manager:StartTournament: error while serializing tournament data: %s", err)
		return
	}

	tm.node.StoreKey(GetTournamentKey(tournament.Id()), data, chord.REPLICATION_FACTOR, chord.UPDATE)
	go NewTournamentRunner(tm, tournament).Resume()
}

func (tm *TournamentManager) SaveTournamentAsJson(id string, json string) {

}

func (tm *TournamentManager) UpdateTournament(tournament models.Tournament) {
	tm.lock.Lock()
	tm.Tournaments[tournament.Id()] = tournament
	tm.lock.Unlock()
	tm.repo.Update(context.Background(), DumpTournament(tournament))

	json, err := TournamentToJson(DumpTournament(tournament))
	if err != nil {
		log.Printf("manager:StartTournament: error while serializing tournament data: %s", err)
		return
	}

	tm.node.StoreKey(GetTournamentKey(tournament.Id()), json, chord.REPLICATION_FACTOR, chord.UPDATE)
}

func (tm *TournamentManager) GetStatus(id string) (int, error) {
	tournament, ok := tm.Tournaments[id]

	if !ok {
		return -1, fmt.Errorf("tournament not found")
	}

	return int(tournament.Status()), nil
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
		tm.lock.Lock()
		tm.Tournaments[tournamentId] = tournament
		tm.lock.Unlock()

		pbTournament := DumpTournament(tournament)

		tm.repo.Update(context.Background(), pbTournament)
		json, err := TournamentToJson(pbTournament)

		if err != nil {
			log.Printf("error serializing tournament")
		}

		tm.node.StoreKey(GetTournamentKey(tournamentId), json, chord.REPLICATION_FACTOR, chord.UPDATE)

		log.Default().Printf("tournamentManager:Notify: tournament %s winner is %s\n", tournament.Id(), tournament.Winner().Id())
		return nil
	}

	if key == "status" {
		tournament.SetStatus(value.(pb.TournamentStatus))
		tm.lock.Lock()
		tm.Tournaments[tournamentId] = tournament
		tm.lock.Unlock()

		pbTournament := DumpTournament(tournament)

		tm.repo.Update(context.Background(), pbTournament)
		json, err := TournamentToJson(pbTournament)

		if err != nil {
			log.Printf("error serializing tournament")
		}

		tm.node.StoreKey(GetTournamentKey(tournamentId), json, chord.REPLICATION_FACTOR, chord.UPDATE)

		log.Default().Printf("tournamentManager:Notify: tournament %s status is %s\n", tournament.Id(), tournament.Status())
	}

	return nil
}
