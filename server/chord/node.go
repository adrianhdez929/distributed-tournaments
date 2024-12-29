package chord

type ChordNode struct {
	client *ChordClient
	server *ChordServer
}

func NewChordNode(ip string) *ChordNode {
	server := NewChordServer(ip, 8080, 10)
	client := NewChordClient(server)

	return &ChordNode{client: client, server: server}
}
