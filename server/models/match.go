package models

import (
	"math/rand"
	"shared/interfaces"
)

type Match interface {
	// Id() string
	Play()
	Next() Match
	Winner() interfaces.Player
	SetPlayer(interfaces.Player)
	Players() [2]interfaces.Player
}

type MatchData struct {
	// id      int
	playerA interfaces.Player
	playerB interfaces.Player
	winner  interfaces.Player
	next    Match
}

func NewMatchData(
	playerA interfaces.Player,
	playerB interfaces.Player,
	next Match,
) *MatchData {
	return &MatchData{playerA, playerB, nil, next}
}

// func (m *MatchData) Id() string {
// 	return fmt.Sprintf("%d", m.id)
// }

func (m *MatchData) Play() {
	m.playerA.Move()
	m.playerB.Move()
	// TODO: improve this random player winner
	r := rand.Intn(2)
	if r == 0 {
		m.winner = m.playerA
		return
	}
	m.winner = m.playerB
}

func (m *MatchData) Next() Match {
	return m.next
}

func (m *MatchData) Winner() interfaces.Player {
	return m.winner
}

func (m *MatchData) Players() [2]interfaces.Player {
	return [2]interfaces.Player{m.playerA, m.playerB}
}

func (m *MatchData) SetPlayer(player interfaces.Player) {
	if m.playerA == nil {
		m.playerA = player
		return
	}

	m.playerB = player
	// TODO: fix bug if more than 2 calls are made
}
