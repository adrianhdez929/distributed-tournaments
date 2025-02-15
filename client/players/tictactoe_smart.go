package players

import (
	"fmt"
	"shared/interfaces"
)

func NewSmartPlayer(id int) interfaces.Player {
	return &TicTacToeGreedyPlayer{id}
}

type TicTacToeSmartPlayer struct {
	id int
}

func (p *TicTacToeSmartPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *TicTacToeSmartPlayer) Move(state interfaces.GameState) interfaces.Move {
	

	return interfaces.Move{X: 2, Y: 3}
}

