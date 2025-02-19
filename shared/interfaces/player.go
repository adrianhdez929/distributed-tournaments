package interfaces

type Player interface {
	Id() string
	Move(GameState) Move
	Name() string
}
