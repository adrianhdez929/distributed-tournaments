package models

import (
	"shared/interfaces"
	"math"
	pb "shared/grpc"
)

type Tournament interface {
	Id() string
	TotalRounds() int
	Players() []interfaces.Player
	Winner() interfaces.Player
	Matches() []Match
	SetMatch(matchId int, match Match)
	Game() string
	SetWinner(winner interfaces.Player)
	SetStatus(status pb.TournamentStatus)
	Status() pb.TournamentStatus
	State() map[string]interface{}
	SetState(key string, value interface{})
	LoadState()
	PendingMatches() map[int]bool
}

var initialState = map[string]interface{}{
	"player_wins": make(map[string]int32),
	"winner":      "",
}

type TournamentData struct {
	id             string
	players        []interfaces.Player
	matches        []Match
	matchCount	   int
	rounds         int
	status         pb.TournamentStatus
	state          map[string]interface{}
	winner         interfaces.Player
	game           string
	pendingMatches map[int]bool	
}


func NewTournamentData(id string, players []interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game) *TournamentData {
	matchC := len(players)-1
	totalRounds := 0
	if len(players)-1 == 0 {
		totalRounds = 0 
	}else{
		totalRounds = int(math.Log2(float64(len(players)-1)))+1
	}
	matches := make([]Match, matchC)
	for i := 0; i < matchC; i++ {
		matches[i] = NewMatchData(gameFactory, i, id)
	}
	cont:=0
	pending := map[int]bool{}
	for i:=0 ;i< matchC; i++{
		if (i+1)*2>matchC{
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if (i+1)*2+1>matchC{
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if matches[i].SetGameIfReady(){
			pending[i]=true
		}		
	}

	return &TournamentData{
		id:             id,
		players:        players,
		matches:        matches,
		matchCount:     matchC,
		rounds:         totalRounds,
		status:         pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
		state:          initialState,
		winner:         nil,
		game:           gameFactory([]interfaces.Player{}).Name(),
		pendingMatches: pending,
	}
}

func (t *TournamentData) Id() string {
	return t.id
}

func (t *TournamentData) Status() pb.TournamentStatus {
	return t.status
}

func (t *TournamentData) SetStatus(status pb.TournamentStatus) {
	t.status = status
}

func (t *TournamentData) Players() []interfaces.Player {
	return t.players
}

func (t *TournamentData) Winner() interfaces.Player {
	return t.winner
}

func (t *TournamentData) TotalRounds() int {
	return t.rounds
}

func (t *TournamentData) Matches() []Match {
	return t.matches[:]
}

func (t *TournamentData) SetMatch(matchId int, match Match) {
	t.matches[matchId] = match
}

func (t *TournamentData) Game() string {
	return t.game
}

func (t *TournamentData) State() map[string]interface{} {
	return t.state
}

func (t *TournamentData) SetState(key string, value interface{}) {
	t.state[key] = value
}

func (t *TournamentData) SetWinner(winner interfaces.Player) {
	t.winner = winner
	t.state["winner"] = winner
}

//funcion para cargar el estado del torneo desde un .json		
func (t *TournamentData) LoadState() {
}

func (t *TournamentData) PendingMatches() map[int]bool {
	return t.pendingMatches
}