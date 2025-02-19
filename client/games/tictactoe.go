package games

import (
	"shared/interfaces"
)

type TicTacToe struct {
	players []interfaces.Player
	status  interfaces.GameStatus
	board   [3][3]int
	winner  interfaces.Player
}

func NewTicTacToe(players []interfaces.Player) interfaces.Game {
	board := [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	return &TicTacToe{players: players, status: interfaces.NotStarted, board: board}
}

func (g *TicTacToe) Play() {
	g.status = interfaces.Running
	lineMade := false

	for !lineMade {
		for i, player := range g.players {
			move := player.Move(g.State())
			g.board[move.X][move.Y] = i
			g.checkLine(g.board, [2]int{move.X, move.Y})
		}
	}
}

func (g *TicTacToe) Winner() interfaces.Player {
	return g.winner
}

func (g *TicTacToe) Status() interfaces.GameStatus {
	return g.status
}

func (g *TicTacToe) State() interfaces.GameState {
	state := make(map[string]interface{})
	state["board"] = g.board
	return state
}

func (g *TicTacToe) Name() string {
	return "TicTacToe"
}

func (g *TicTacToe) checkLine(matrix [3][3]int, coord [2]int) bool {
	directions := [4][2]int{{1, 0}, {0, 1}, {0, 0}, {1, 1}}
	for _, direction := range directions {
		checkCount := 1

		if matrix[abs((coord[0]+direction[0])%2)][abs((coord[1]+direction[1])%2)] == matrix[coord[0]][coord[1]] {
			checkCount++
		}

		if matrix[abs((coord[0]-direction[0])%2)][abs((coord[1]-direction[1])%2)] == matrix[coord[0]][coord[1]] {
			checkCount++
		}

		if checkCount == 3 {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
