package chord

type ChordClient struct {
	server ChordServer
}

func NewChordClient(server ChordServer) *ChordClient {
	return &ChordClient{server: server}
}
