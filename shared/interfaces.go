package shared

type Player interface {
	Id() string
	Move()
}

type Game interface {
}
