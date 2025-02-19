package models

import (
	"math"
	"shared/interfaces"

	pb "shared/grpc"
)

type Tournament interface {
	Id() string
	CurrentRound() int
	TotalRounds() int
	Players() []interfaces.Player
	Winner() interfaces.Player
	Matches() map[string]Match
	SetMatch(matchId string, match Match)
	InitialMatches() []Match
	Game() string
	// AddMatches(matches []Match)
	SetWinner(winner interfaces.Player)
	SetStatus(status pb.TournamentStatus)
	Status() pb.TournamentStatus
	State() map[string]interface{}
	SetState(key string, value interface{})
}

var initialState = map[string]interface{}{
	"player_wins": make(map[string]int32),
	"winner":      "",
}

type TournamentData struct {
	id             string
	players        []interfaces.Player
	matches        map[string]Match
	initialMatches []Match
	currentRound   int
	rounds         int
	status         pb.TournamentStatus
	winner         interfaces.Player
	state          map[string]interface{}
	game           string
}

func NewTournamentData(id string, players []interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game) *TournamentData {
	initialMatchCount := float64(len(players) / 2)

	totalRounds := int(math.Log2(initialMatchCount)) + 1
	matches := createTournament(totalRounds, gameFactory)

	for _, v := range matches {
		v.SetPlayer(players[0])
		v.SetPlayer(players[1])
		players = players[2:]
	}

	return &TournamentData{
		id:             id,
		players:        players,
		initialMatches: matches,
		matches:        make(map[string]Match),
		currentRound:   1,
		rounds:         totalRounds,
		status:         pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
		state:          initialState,
		game:           gameFactory([]interfaces.Player{}).Name(),
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

func (t *TournamentData) CurrentRound() int {
	return t.currentRound
}

func (t *TournamentData) TotalRounds() int {
	return t.rounds
}

func (t *TournamentData) Matches() map[string]Match {
	return t.matches
}

func (t *TournamentData) SetMatch(matchId string, match Match) {
	t.matches[matchId] = match
}

func (t *TournamentData) InitialMatches() []Match {
	return t.initialMatches
}

func (t *TournamentData) Game() string {
	return t.game
}

// func (t *TournamentData) AddMatches(matches []Match) {
// 	t.matches = append(t.matches, matches...)
// }

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

func createTournament(rounds int, gameFactory func([]interfaces.Player) interfaces.Game) []Match {
	finalMatch := NewMatchData(gameFactory, nil, nil)
	currentRound := []Match{finalMatch}

	for i := 0; i < rounds-1; i++ {
		newRound := make([]Match, 0, 2*len(currentRound))

		for _, v := range currentRound {
			submatches := createSubMatches(v, gameFactory)
			newRound = append(newRound, submatches[0], submatches[1])
		}

		currentRound = newRound
	}

	return currentRound
}

func createSubMatches(match Match, gameFactory func([]interfaces.Player) interfaces.Game) [2]Match {
	return [2]Match{
		NewMatchData(gameFactory, nil, match),
		NewMatchData(gameFactory, nil, match),
	}
}
