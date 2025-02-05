package chord

type ChordNode struct {
	client *ChordClient
	server *ChordServer
}

func NewChordNode(ip string, port int) *ChordNode {
	server := NewChordServer(ip, port, 160)
	client := NewChordClient(server)

	return &ChordNode{client: client, server: server}
}
