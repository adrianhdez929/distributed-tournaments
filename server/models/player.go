package models

import (
	"encoding/json"
	"fmt"
)

type PlayerData struct {
	Id_ int `json:"id"`
}

func (p *PlayerData) FromJson(jsonData string) *PlayerData {
	data := &PlayerData{}

	err := json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil
	}
	return data
}

func (p *PlayerData) ToJson() string {
	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(data)
}

func NewPlayerData(id int) *PlayerData {
	return &PlayerData{Id_: id}
}

func (p *PlayerData) Id() string {
	return fmt.Sprintf("%d", p.Id_)
}

func (p *PlayerData) Move() {}
