package models

import (
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
	Status()	int
	SetStatus(int)
	Start() interfaces.Player
}

type MatchData struct {
	id          int
	tournamentId string
	gameFactory func([]interfaces.Player) interfaces.Game
	game        interfaces.Game
	players     []interfaces.Player
	matchStatus int
	winner      interfaces.Player
}

func NewMatchData(
	gameFactory func([]interfaces.Player) interfaces.Game,
	matchId int,
	tournamentID string,
) *MatchData {
	return &MatchData{
		id:				matchId,
		tournamentId: 	tournamentID,
		gameFactory:  	gameFactory,
		game:			nil,
		players:		[]interfaces.Player{},
		matchStatus: 	0,
		winner: 		nil,
	}
}

func (m *MatchData) Id() int{
	return m.id
}

func (m *MatchData) TournamentId() string {
	return m.tournamentId
}

func (m *MatchData) Game() interfaces.Game {
	return m.game
}

func (m *MatchData) SetGameIfReady() bool {
	if(len(m.players) < 2){
		return false
	}
	m.matchStatus=1
	m.game=m.gameFactory(m.players)
	return true
}

func (m *MatchData) Winner() interfaces.Player {
	return m.winner
}

func (m *MatchData) Players() []interfaces.Player {
	return m.players
}

func (m *MatchData) AddPlayer(player interfaces.Player) {
	m.players = append(m.players, player)
}

func (m *MatchData) SetWinner(player interfaces.Player) {
	m.winner = player
}

//funcion para obtener el estado de la partida
func (m *MatchData) Status() int {
	return m.matchStatus
}
//funcion para settear el estado de la partida
func (m *MatchData) SetStatus(status int) {	
	m.matchStatus = status	
}

//funcion para cargar los datos de la partida desde un .json
func (m *MatchData) LoadState() {
}

func (m *MatchData) Start() interfaces.Player {
	m.game.Play()
	return m.game.Winner()
}