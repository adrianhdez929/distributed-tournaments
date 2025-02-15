package players

import (
	"fmt"
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
	directions := [4][2]int{{1, 1}, {0, 1}, {1, -1}, {1, 0}}
	Ans:= [6]int{-1,-1,-1,-1,-1,-1}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				candidate:=[6]int{-1,-1,-1,-1,i,j}
				for dir:=0; dir<4; dir++{
					val:=0
					oc:=0
					for mov:=-2;mov<=2;mov++{
						if !check(i+directions[dir][0]*mov,j+directions[dir][1]*mov){
							continue
						}
						if board[i+directions[dir][0]*mov][j+directions[dir][1]*mov]==0{
							val++
							continue
						}
						if board[i+directions[dir][0]*mov][j+directions[dir][1]*mov]!=p.id{
							break;
						}
						val++
						oc++
					}
					if val==3{
						candidate[dir]=oc
					}
				}
				candidate = parcial_sort(candidate)
				Ans= Max(Ans,candidate)	
			}
		}
	}

	return interfaces.Move{X: Ans[4], Y: Ans[5]}
}

func check( x int, y int) bool{
	return x>=0 && x<3 && y>=0 && y<3
}

func parcial_sort(candidate [6]int) [6]int{
	for i:=0;i<3;i++{
		for j:=0;j<3;j++{
			if candidate[j]<candidate[j+1]{
				candidate[j],candidate[j+1] = candidate[j+1],candidate[j]
			}
		}
	}
	return candidate
}

func Max(A [6]int, B [6]int) [6]int{
	for i:=0 ; i<6 ;i++{
		if(A[i]<B[i]){
			return B
		}
		if(A[i]>B[i]){
			return A;
		}
	}
	return A;
}