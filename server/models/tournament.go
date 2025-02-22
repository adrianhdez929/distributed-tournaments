package models

import (
	"encoding/json"
	"math"
	pb "shared/grpc"
	"shared/interfaces"
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
	Id_             string                 `json:"id"`
	Players_        []interfaces.Player    `json:"players"`
	Matches_        []Match                `json:"matches"`
	MatchCount_     int                    `json:"match_count"`
	Rounds_         int                    `json:"rounds"`
	Status_         pb.TournamentStatus    `json:"status"`
	State_          map[string]interface{} `json:"state"`
	Winner_         interfaces.Player      `json:"winner"`
	Game_           string                 `json:"game"`
	PendingMatches_ map[int]bool           `json:"pending_matches"`
}

func NewTournamentData(id string, players []interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game) *TournamentData {
	matchC := len(players) - 1
	totalRounds := 0
	if len(players)-1 == 0 {
		totalRounds = 0
	} else {
		totalRounds = int(math.Log2(float64(len(players)-1))) + 1
	}
	matches := make([]Match, matchC)
	for i := 0; i < matchC; i++ {
		matches[i] = NewMatchData(gameFactory, i, id)
	}
	cont := 0
	pending := map[int]bool{}
	for i := 0; i < matchC; i++ {
		if (i+1)*2 > matchC {
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if (i+1)*2+1 > matchC {
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if matches[i].SetGameIfReady() {
			pending[i] = true
		}
	}

	return &TournamentData{
		Id_:             id,
		Players_:        players,
		Matches_:        matches,
		MatchCount_:     matchC,
		Rounds_:         totalRounds,
		Status_:         pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
		State_:          initialState,
		Winner_:         nil,
		Game_:           gameFactory([]interfaces.Player{}).Name(),
		PendingMatches_: pending,
	}
}

func (t *TournamentData) FromJson(jsonData string) *TournamentData {
	var data *TournamentData

	err := json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil
	}
	return data
}

func (t *TournamentData) ToJson() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(data)
}

func (t *TournamentData) Id() string {
	return t.Id_
}

func (t *TournamentData) Status() pb.TournamentStatus {
	return t.Status_
}

func (t *TournamentData) SetStatus(status pb.TournamentStatus) {
	t.Status_ = status
}

func (t *TournamentData) Players() []interfaces.Player {
	return t.Players_
}

func (t *TournamentData) Winner() interfaces.Player {
	return t.Winner_
}

func (t *TournamentData) TotalRounds() int {
	return t.Rounds_
}

func (t *TournamentData) Matches() []Match {
	return t.Matches_[:]
}

func (t *TournamentData) SetMatch(matchId int, match Match) {
	t.Matches_[matchId] = match
}

func (t *TournamentData) Game() string {
	return t.Game_
}

func (t *TournamentData) State() map[string]interface{} {
	return t.State_
}

func (t *TournamentData) SetState(key string, value interface{}) {
	t.State_[key] = value
}

func (t *TournamentData) SetWinner(winner interfaces.Player) {
	t.Winner_ = winner
	t.State_["winner"] = winner
}

// funcion para cargar el estado del torneo desde un .json
func (t *TournamentData) LoadState() {
}

func (t *TournamentData) PendingMatches() map[int]bool {
	return t.PendingMatches_
}
