package players

import (
	"fmt"
	"math/rand"
	"shared/interfaces"
)

func NewHexRandomPlayer(id int) interfaces.Player {
	return &TicTacToeRandomPlayer{id}
}

type HexRandomPlayer struct {
	id int
}

func (p *HexRandomPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *HexRandomPlayer) Move(state interfaces.GameState) interfaces.Move {
	board := state["board"].([19][19]int)

	freeCells := p.getFreeCells(board)
	index := rand.Int() % len(freeCells)

	return freeCells[index]
}

func (p *HexRandomPlayer)getFreeCells(board [19][19]int) []interfaces.Move {
	freeCells := []interfaces.Move{}
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if board[i][j] == 0 {
				freeCells = append(freeCells, interfaces.Move{X: i, Y: j})
			}
		}
	}
	return freeCells
}
