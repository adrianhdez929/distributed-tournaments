package players

import (
	"fmt"
	"math/rand"
	"shared/interfaces"
)

func NewRandomPlayer(id int) interfaces.Player {
	return &TicTacToeRandomPlayer{id}
}

type TicTacToeRandomPlayer struct {
	id int
}

func (p *TicTacToeRandomPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *TicTacToeRandomPlayer) Move(state interfaces.GameState) interfaces.Move {
	board := state["board"].([3][3]int)
	valid := false
	for !valid {
		x := rand.Int() % 2
		y := rand.Int() % 2

		if board[x][y] == 0 {
			valid = true
			return interfaces.Move{X: x, Y: y}
		}
	}

	return interfaces.Move{X: 0, Y: 0}
}

func (p *TicTacToeRandomPlayer) Name() string {
	return "TicTacToeRandomPlayer"
}
