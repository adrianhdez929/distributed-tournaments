package main

import (
	"fmt"
	"tournaments/server/models"
)

func NewGreedyPlayer(id int) models.Player {
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
