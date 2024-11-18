package models

import "fmt"

type PlayerData struct {
	id int
}

func NewPlayerData(id int) *PlayerData {
	return &PlayerData{id}
}

func (p *PlayerData) Id() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *PlayerData) Move() {}
