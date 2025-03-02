package games

import (
	"shared/interfaces"
)

var gameConstructor = map[string]func([]interfaces.Player) interfaces.Game{
	"TicTacToe": NewTicTacToe,
}

func NewGameFactory(game string, players []interfaces.Player) func([]interfaces.Player) interfaces.Game {
	return gameConstructor[game]
}
