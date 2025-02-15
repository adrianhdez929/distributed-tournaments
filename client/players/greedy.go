package players

import (
	"fmt"
	"math/rand"
	"shared/interfaces"
)

func NewGreedyPlayer(id int) interfaces.Player {
	return &TicTacToeGreedyPlayer{id}
}

type TicTacToeGreedyPlayer struct {
	id int
}

func (p *TicTacToeGreedyPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *TicTacToeGreedyPlayer) Move(state interfaces.GameState) interfaces.Move {
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
