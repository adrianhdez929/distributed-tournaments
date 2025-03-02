package players

import (
	"shared/interfaces"
)

var playerConstructor = map[string]func(id string) interfaces.Player{
	"TicTacToeRandomPlayer": NewRandomPlayer,
}

func NewPlayerFactory(agentName string, name string) interfaces.Player {
	return playerConstructor[agentName](name)
}
