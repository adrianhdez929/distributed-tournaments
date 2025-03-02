package models

import (
	"encoding/json"
	"shared/interfaces"
)

type Match interface {
	Id() int
	TournamentId() string
	Game() interfaces.Game
	Winner() interfaces.Player
	AddPlayer(interfaces.Player)
	SetWinner(interfaces.Player)
	Players() []interfaces.Player
	SetGameIfReady() bool
	Status() int
	SetStatus(int)
	Start() interfaces.Player
	ToJson() string
}

type MatchData struct {
	Id_           int    `json:"id"`
	TournamentId_ string `json:"tournament_id"`
	gameFactory   func([]interfaces.Player) interfaces.Game
	Game_         interfaces.Game     `json:"game"`
	Players_      []interfaces.Player `json:"players"`
	MatchStatus_  int                 `json:"status"`
	Winner_       interfaces.Player   `json:"winner"`
}

func NewMatchData(
	gameFactory func([]interfaces.Player) interfaces.Game,
	matchId int,
	tournamentID string,
) *MatchData {
	return &MatchData{
		Id_:           matchId,
		TournamentId_: tournamentID,
		gameFactory:   gameFactory,
		Game_:         nil,
		Players_:      []interfaces.Player{},
		MatchStatus_:  0,
		Winner_:       nil,
	}
}

func (m MatchData) FromJson(jsonData string) *MatchData {
	match := &MatchData{}

	err := json.Unmarshal([]byte(jsonData), match)
	if err != nil {
		return nil
	}
	return match
}

func (m MatchData) ToJson() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m *MatchData) Id() int {
	return m.Id_
}

func (m *MatchData) TournamentId() string {
	return m.TournamentId_
}

func (m *MatchData) Game() interfaces.Game {
	return m.Game_
}

func (m *MatchData) SetGameIfReady() bool {
	if len(m.Players_) < 2 {
		return false
	}
	m.MatchStatus_ = 1
	m.Game_ = m.gameFactory(m.Players_)
	return true
}

func (m *MatchData) Winner() interfaces.Player {
	return m.Winner_
}

func (m *MatchData) Players() []interfaces.Player {
	return m.Players_
}

func (m *MatchData) AddPlayer(player interfaces.Player) {
	m.Players_ = append(m.Players_, player)
}

func (m *MatchData) SetWinner(player interfaces.Player) {
	m.Winner_ = player
}

// funcion para obtener el estado de la partida
func (m *MatchData) Status() int {
	return m.MatchStatus_
}

// funcion para settear el estado de la partida
func (m *MatchData) SetStatus(status int) {
	m.MatchStatus_ = status
}

func (m *MatchData) Start() interfaces.Player {
	m.Game_.Play()
	return m.Game_.Winner()
}
