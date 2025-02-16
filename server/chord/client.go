package chord

type ChordClient struct {
	server *ChordServer
}

func NewChordClient(server *ChordServer) *ChordClient {
	return &ChordClient{server: server}
}

func (c *ChordClient) GetSuccessor() ChordNodeReference {
	return c.server.Successor()
}
