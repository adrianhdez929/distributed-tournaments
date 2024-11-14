package models

// "time"

type Match interface {
	// Id() string
	Play()
	Next() Match
	Winner() Player
	SetPlayer(Player)
	Players() [2]Player
}

type MatchData struct {
	// id      int
	playerA Player
	playerB Player
	winner  Player
	next    Match
}

func NewMatchData(playerA Player, playerB Player, next Match) *MatchData {
	return &MatchData{playerA, playerB, nil, next}
}

// func (m *MatchData) Id() string {
// 	return fmt.Sprintf("%d", m.id)
// }

func (m *MatchData) Play() {
	m.winner = m.playerA
}

func (m *MatchData) Next() Match {
	return m.next
}

func (m *MatchData) Winner() Player {
	return m.winner
}

func (m *MatchData) Players() [2]Player {
	return [2]Player{m.playerA, m.playerB}
}

func (m *MatchData) SetPlayer(player Player) {
	if m.playerA == nil {
		m.playerA = player
		return
	}

	m.playerB = player
	// TODO: fix bug if more than 2 calls are made
}
