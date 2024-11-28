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

	freeCells := getFreeCells(board)
	index := rand.Int() % len(freeCells)

	return freeCells[index]
}

func getFreeCells(board [3][3]int) []interfaces.Move {
	freeCells := []interfaces.Move{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				freeCells = append(freeCells, interfaces.Move{X: i, Y: j})
			}
		}
	}
	return freeCells
}
