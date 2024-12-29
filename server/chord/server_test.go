package chord_test

import (
	"testing"
	"tournament_server/chord"
)

func prepareChordServer() *chord.ChordServer {
	server := chord.NewChordServer("127.0.0.1", 8080, 10)

	return server
}

func TestChordServer(t *testing.T) {
	server := prepareChordServer()

	t.Run("should start whitout crashing", func(t *testing.T) {
		if server.Successor().Id == 0 {
			t.Errorf("Successor is not valid")
		}
	})

	t.Run("should return self as successor if it is the only node", func(t *testing.T) {
		if server.Successor().Id != server.Id() {
			t.Errorf("Successor is not self reference")
		}
	})

	t.Run("should join another node", func(t *testing.T) {

	})
}
