package players

import (
	"fmt"
)

func NewGreedyPlayer(id int) *GreedyPlayer {
	return &GreedyPlayer{id}
}

type GreedyPlayer struct {
	id int
}

func (p *GreedyPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *GreedyPlayer) Move() {
	fmt.Printf("%s is moving", p.Id())
}
