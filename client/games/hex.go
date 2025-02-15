package games

import (
	"shared/interfaces"
)

type Hex struct {
	players []interfaces.Player
	status  interfaces.GameStatus
	board   [19][19]int
	winner  interfaces.Player
}

func NewHex(players []interfaces.Player) interfaces.Game {
	board := [19][19]int{}
	return &Hex{
		players: players,
		status:  interfaces.NotStarted,
		board:   board,
		winner:  nil,
	}
}

func (g *Hex) Play() {
	// Inicializa el estado del juego
	g.status = interfaces.Running

	// Ciclo principal del juego
	for {
		// Itera sobre los jugadores
		for i, player := range g.players {
			// Obtiene el movimiento del jugador
			move := player.Move(g.State())

			// Verifica si el movimiento es válido
			if !g.isValidHexMove(move) {
				// Si no es válido, ignora el movimiento
				return
			}

			// Actualiza el tablero con el movimiento
			g.board[move.X][move.Y] = i

			// Verifica si hay un ganador
			if g.hasWinner() {
				// Si hay un ganador, actualiza el estado del juego y devuelve el ganador
				g.status = interfaces.Finished
				g.winner = g.players[i]
				return
			}
		}
	}
}

func (g *Hex)isValidHexMove(move interfaces.Move) bool {
	// Verifica si la posición está dentro del tablero
	if move.X < 0 || move.X >= 19 || move.Y < 0 || move.Y >= 19 {
		return false
	}

	// Verifica si la posición está vacía
	if g.board[move.X][move.Y] != 0 {
		return false
	}

	return true
}


func (g *Hex) hasWinner() bool {
    // Verifica si hay un ganador en la dirección horizontal
    mk := [19][19]int{}
	for i := 0; i < 19; i++ {
        if g.board[i][0] == 1 {
            if g.dfs(i, 0, 1, mk) {
				return true
            }
        }
    }
	mk= [19][19]int{}
    // Verifica si hay un ganador en la dirección vertical
    for j := 0; j < 19; j++ {
        if g.board[0][j] == 2 {
            if g.dfs(0, j, 2, mk) {
                return true
            }
        }
    }

    return false
}

func (g *Hex) dfs(x int, y int, color int, mk [19][19]int) bool {
    
	directions := [6][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {-1,1}, {1,-1}}

	if x < 0 || x >= 19 || y < 0 || y >= 19 {
        return false
    }
	if mk[x][y] == 1 {
		return false
	}
    if g.board[x][y] != color {
        return false
    }
	mk[x][y] = 1

	if color == 1 && y == 18 {
		return true
	}
	if color == 2 && x == 18 {
		return true
	}  

	for i:=0; i<6 ;i++{
		if g.dfs(x+directions[i][0],y+directions[i][1],color,mk){
			return true
		}
	}
	return false
}


func (g *Hex) Winner() interfaces.Player {
	return g.winner
}

func (g *Hex) Status() interfaces.GameStatus {
	return g.status
}

func (g *Hex) State() interfaces.GameState {
	state := make(map[string]interface{})
	return state
}
