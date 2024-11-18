package players

import (
	"fmt"
	"shared/interfaces"
)

func NewGreedyPlayer(id int) interfaces.Player {
	return &GreedyPlayer{id}
}

type GreedyPlayer struct {
	id int
}

func (p *GreedyPlayer) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *GreedyPlayer) Move() {
	fmt.Printf("%s is moving\n", p.Id())
}
