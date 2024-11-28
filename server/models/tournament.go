package models

import (
	"math"
	"shared/interfaces"
)

type Tournament interface {
	CurrentRound() int
	TotalRounds() int
	Winner() interfaces.Player
}

type TournamentData struct {
	players      []interfaces.Player
	matches      []Match
	winner       interfaces.Player
	currentRound int
	rounds       int
}

func NewTournamentData(players []interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game) *TournamentData {
	initialMatchCount := float64(len(players) / 2)

	totalRounds := int(math.Log2(initialMatchCount)) + 1
	matches := createTournament(totalRounds, gameFactory)

	for _, v := range matches {
		v.SetPlayer(players[0])
		v.SetPlayer(players[1])
		players = players[2:]
	}

	return &TournamentData{
		players,
		matches,
		nil,
		1,
		totalRounds,
	}
}

func (t *TournamentData) Winner() interfaces.Player {
	if t.winner != nil {
		return t.winner
	}

	currentMatches := t.matches

	for i := 0; i < t.rounds; i++ {
		newMatches := make([]Match, 0, len(currentMatches)/2)

		for _, v := range currentMatches {
			v.Start()
			winner := v.Winner()

			if v.Next() != nil {
				v.Next().SetPlayer(winner)
			} else {
				t.winner = winner
			}

			if len(newMatches) == 0 {
				newMatches = append(newMatches, v.Next())
			} else if newMatches[len(newMatches)-1] != v.Next() {
				newMatches = append(newMatches, v.Next())
			}
		}

		currentMatches = newMatches
	}

	return t.winner
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
