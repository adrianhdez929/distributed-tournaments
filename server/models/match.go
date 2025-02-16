package models

import (
	"fmt"
	"math/rand"
	"shared/interfaces"
)

type Match interface {
	Id() string
	Game() interfaces.Game
	Start() interfaces.Player
	Next() Match
	Winner() interfaces.Player
	SetPlayer(interfaces.Player)
	SetWinner(interfaces.Player)
	Players() []interfaces.Player
}

type MatchData struct {
	id          int
	gameFactory func([]interfaces.Player) interfaces.Game
	game        interfaces.Game
	players     []interfaces.Player
	winner      interfaces.Player
	next        Match
}

func NewMatchData(
	gameFactory func([]interfaces.Player) interfaces.Game,
	players []interfaces.Player,
	next Match,
) *MatchData {
	return &MatchData{
		// TODO: hash id
		rand.Int(),
		gameFactory,
		gameFactory(players),
		players,
		nil,
		next}
}

func (m *MatchData) Id() string {
	return fmt.Sprintf("%d", m.id)
}

func (m *MatchData) Game() interfaces.Game {
	return m.game
}

func (m *MatchData) Start() interfaces.Player {
	m.game = m.gameFactory(m.players)
	m.game.Play()
	m.winner = m.game.Winner()
	fmt.Println("winner is ", m.Winner().Id())

	return m.winner
}

func (m *MatchData) Next() Match {
	return m.next
}

func (m *MatchData) Winner() interfaces.Player {
	return m.winner
}

func (m *MatchData) Players() []interfaces.Player {
	return m.players
}

func (m *MatchData) SetPlayer(player interfaces.Player) {
	m.players = append(m.players, player)
}

func (m *MatchData) SetWinner(player interfaces.Player) {
	m.winner = player
}
