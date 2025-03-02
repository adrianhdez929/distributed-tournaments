package models

import (
	"encoding/json"
	"log"
	"math"
	pb "shared/grpc"
	"shared/interfaces"
	"strconv"
	"tournament_server/games"
	"tournament_server/players"
)

type Tournament interface {
	Id() string
	TotalRounds() int
	Players() []interfaces.Player
	Winner() interfaces.Player
	Matches() []Match
	SetMatch(matchId int, match Match)
	Game() string
	SetWinner(winner interfaces.Player)
	SetStatus(status pb.TournamentStatus)
	Status() pb.TournamentStatus
	State() map[string]interface{}
	SetState(key string, value interface{})
	PendingMatches() map[int]bool
	ToJson() string
}

var initialState = map[string]interface{}{
	"player_wins": make(map[string]int32),
	"winner":      "",
}

type TournamentData struct {
	Id_             string                 `json:"id"`
	Players_        []interfaces.Player    `json:"players"`
	Matches_        []Match                `json:"matches"`
	MatchCount_     int                    `json:"match_count"`
	Rounds_         int                    `json:"rounds"`
	Status_         pb.TournamentStatus    `json:"status"`
	State_          map[string]interface{} `json:"state"`
	Winner_         interfaces.Player      `json:"winner"`
	Game_           string                 `json:"game"`
	PendingMatches_ map[int]bool           `json:"pending_matches"`
}

func NewTournamentData(id string, players []interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game) *TournamentData {
	matchC := len(players) - 1
	totalRounds := 0
	if len(players)-1 == 0 {
		totalRounds = 0
	} else {
		totalRounds = int(math.Log2(float64(len(players)-1))) + 1
	}
	matches := make([]Match, matchC)
	for i := 0; i < matchC; i++ {
		matches[i] = NewMatchData(gameFactory, i, id)
	}
	cont := 0
	pending := map[int]bool{}
	for i := 0; i < matchC; i++ {
		if (i+1)*2 > matchC {
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if (i+1)*2+1 > matchC {
			matches[i].AddPlayer(players[cont])
			cont++
		}
		if matches[i].SetGameIfReady() {
			pending[i] = true
		}
	}

	return &TournamentData{
		Id_:             id,
		Players_:        players,
		Matches_:        matches,
		MatchCount_:     matchC,
		Rounds_:         totalRounds,
		Status_:         pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
		State_:          initialState,
		Winner_:         nil,
		Game_:           gameFactory([]interfaces.Player{}).Name(),
		PendingMatches_: pending,
	}
}

func (t TournamentData) FromJson(jsonData string) *TournamentData {
	// First unmarshal into a map to handle interface conversion
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &jsonMap)
	if err != nil {
		log.Printf("Error unmarshaling tournament: %s", err)
		return nil
	}

	// Create tournament data
	data := &TournamentData{}

	// Convert basic fields
	data.Id_ = jsonMap["id"].(string)
	data.Game_ = jsonMap["game"].(string)

	if jsonMap["status"] != nil {
		data.Status_ = pb.TournamentStatus(pb.TournamentStatus_value[jsonMap["status"].(string)])
	} else {
		data.Status_ = pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED
	}
	if jsonMap["state"] != nil {
		data.State_ = jsonMap["state"].(map[string]interface{})
	} else {
		data.State_ = make(map[string]interface{})
	}
	data.PendingMatches_ = make(map[int]bool)

	// Convert players
	playersJson := jsonMap["players"].([]interface{})
	data.Players_ = make([]interfaces.Player, len(playersJson))
	for i, p := range playersJson {
		playerMap := p.(map[string]interface{})
		player := players.NewPlayerFactory(playerMap["agentName"].(string), playerMap["name"].(string))
		data.Players_[i] = player
	}

	// Convert matches if they exist
	if matchesJson, ok := jsonMap["matches"].([]interface{}); ok {
		data.Matches_ = make([]Match, len(matchesJson))
		data.PendingMatches_ = make(map[int]bool, len(matchesJson))

		for i, m := range matchesJson {
			matchMap := m.(map[string]interface{})
			mId, _ := matchMap["id"].(string)
			id, _ := strconv.Atoi(mId)
			match := NewMatchData(
				games.NewGameFactory(jsonMap["game"].(string), []interfaces.Player{}),
				id,
				data.Id_,
			)

			// Add players to match if they exist
			if p1, ok := matchMap["player1"].(map[string]interface{}); ok {
				player1 := players.NewPlayerFactory(p1["agentName"].(string), p1["name"].(string))
				match.AddPlayer(player1)
			}
			if p2, ok := matchMap["player2"].(map[string]interface{}); ok {
				player2 := players.NewPlayerFactory(p2["agentName"].(string), p2["name"].(string))
				match.AddPlayer(player2)
			}

			data.Matches_[i] = match
			if match.SetGameIfReady() {
				data.PendingMatches_[i] = true
			}
		}
	}

	return data
}

func (t *TournamentData) ToJson() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(data)
}

func (t *TournamentData) Id() string {
	return t.Id_
}

func (t *TournamentData) Status() pb.TournamentStatus {
	return t.Status_
}

func (t *TournamentData) SetStatus(status pb.TournamentStatus) {
	t.Status_ = status
}

func (t *TournamentData) Players() []interfaces.Player {
	return t.Players_
}

func (t *TournamentData) Winner() interfaces.Player {
	return t.Winner_
}

func (t *TournamentData) TotalRounds() int {
	return t.Rounds_
}

func (t *TournamentData) Matches() []Match {
	return t.Matches_[:]
}

func (t *TournamentData) SetMatch(matchId int, match Match) {
	t.Matches_[matchId] = match
}

func (t *TournamentData) Game() string {
	return t.Game_
}

func (t *TournamentData) State() map[string]interface{} {
	return t.State_
}

func (t *TournamentData) SetState(key string, value interface{}) {
	t.State_[key] = value
}

func (t *TournamentData) SetWinner(winner interfaces.Player) {
	t.Winner_ = winner
	t.State_["winner"] = winner
}

func (t *TournamentData) PendingMatches() map[int]bool {
	return t.PendingMatches_
}
