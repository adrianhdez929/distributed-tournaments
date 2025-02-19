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
	board := state["board"].([3][3]int)
	move,_:=p.BestMove(board,p.id)
	return move
}

func (p *TicTacToeSmartPlayer) BestMove(board [3][3]int,id int) (interfaces.Move,int){
	other:=-1
	for i:=0;i<3;i++{
		for j:=0;j<3;j++{
			if board[i][j]!=0 && board[i][j]!=id{
				other=board[i][j]
				break
			}
		}
	}
	ans1:= interfaces.Move{X: -1, Y: -1,Value: id}
	val:=-1
	for i:=0;i<3;i++{
		for j:=0;j<3;j++{
			if board[i][j]==0{
				board[i][j]=id
				_,value :=p.BestMove(board,other)
				if value==-1{
					return interfaces.Move{X: i, Y: j,Value: id},1
				}
				if value==0{
					ans1=interfaces.Move{X: i, Y: j,Value: id}
					val=0
				}
				if val==-1{
					ans1=interfaces.Move{X: i, Y: j,Value: id}
				}
			}
		}
	}
	if ans1.X==-1{
		return ans1,0
	}
	return ans1,val
	
}

