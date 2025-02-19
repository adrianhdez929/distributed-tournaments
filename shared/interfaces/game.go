package interfaces

type GameStatus string
type GameState map[string]interface{}

const (
	NotStarted GameStatus = "NOT_STARTED"
	Running    GameStatus = "RUNNING"
	Finished   GameStatus = "FINISHED"
)

type Game interface {
	Status() GameStatus
	State() GameState
	Play()
	Winner() Player
	Name() string
}
