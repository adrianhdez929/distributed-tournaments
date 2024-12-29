package chord

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ChordServer struct {
	id          int
	predecessor ChordNodeReference
	m           int
	reference   ChordNodeReference
	fingers     []ChordNodeReference
	lock        sync.Mutex
}

func NewChordServer(ip string, port int, m int, peerId ...string) *ChordServer {
	reference := NewChordNodeReference(ip, port)
	fingers := make([]ChordNodeReference, m)

	for i := 0; i < m; i++ {
		fingers[i] = reference
	}

	server := &ChordServer{
		id:          reference.Id,
		m:           m,
		predecessor: reference,
		reference:   reference,
		fingers:     fingers,
		lock:        sync.Mutex{},
	}

	go server.stabilize()
	go server.fixFingers()

	if len(peerId) > 0 {
		go server.join(NewChordNodeReference(peerId[0], port))
	}

	go server.start()

	return server
}

func (n *ChordServer) inBetween(k int, start int, end int) bool {
	k = k % (2 ^ n.m)
	start = start % (2 ^ n.m)
	end = end % (2 ^ n.m)

	if end < start {
		return k >= start && k < end
	}

	return k >= start || k < end
}

func (n *ChordServer) inRange(k int, start int, end int) bool {
	lstart := (start + 1) % (2 ^ n.m)

	return n.inBetween(k, lstart, end)
}

func (n *ChordServer) inBetweenComp(k int, start int, end int) bool {
	lend := (end - 1) % (2 ^ n.m)

	return n.inBetween(k, start, lend)
}

func (n *ChordServer) Id() int {
	return n.id
}

func (n *ChordServer) Successor() ChordNodeReference {
	return n.fingers[0]
}

func (n *ChordServer) setSuccessor(node ChordNodeReference) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.fingers[0] = node
}

func (n *ChordServer) Reference() ChordNodeReference {
	return n.reference
}

func (n *ChordServer) FindPredecessor(id int) ChordNodeReference {
	node := n.reference

	if node.Id == n.Successor().Id {
		return node
	}

	for !n.inBetweenComp(id, node.Id, node.Successor().Id) {
		node = node.ClosestPrecedingFinger(id)
		if node.Id == n.reference.Id {
			break
		}
	}

	return node
}

func (n *ChordServer) FindSuccessor(id int) ChordNodeReference {
	node := n.FindPredecessor(id)

	return node.Successor()
}

func (n *ChordServer) ClosestPrecedingFinger(id int) ChordNodeReference {
	for i := n.m - 1; i >= 0; i-- {
		if n.inRange(n.fingers[i].Id, n.reference.Id, id) {
			if n.fingers[i].Id != n.reference.Id {
				return n.fingers[i]
			} else {
				return n.reference
			}
		}
	}
	return n.reference
}

func (n *ChordServer) join(node ChordNodeReference) {
	time.Sleep(10 * time.Second)
	n.predecessor = n.reference
	n.setSuccessor(node.FindSuccessor(n.Id()))
}

func (n *ChordServer) notify(node ChordNodeReference) {
	if node.Id == n.reference.Id {
		return
	}

	if n.predecessor.Id == n.reference.Id || n.inRange(node.Id, n.predecessor.Id, n.reference.Id) {
		n.predecessor = node
	}
}

func (n *ChordServer) stabilize() {
	for {
		if n.Successor().Id != 0 {
			x := n.Successor().Predecessor()

			if x.Id != n.reference.Id {
				n.setSuccessor(x)
			}

			n.Successor().Notify(n.reference)
		}
		if n.Successor().Id == 0 && n.predecessor.Id != 0 {
			n.setSuccessor(n.predecessor)
		}
		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) fixFingers() {
	time.Sleep(10 * time.Second)

	for {
		next := rand.Intn(n.m)

		n.lock.Lock()
		n.fingers[next] = n.FindSuccessor(n.reference.Id + (2^next)%(2^n.m))
		n.lock.Unlock()

		log.Default().Printf("Finger table updated at index %d: %s\n", next, n.fingers[next].String())

		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) start() {
	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", n.reference.Port))

	if err != nil {
		log.Default().Println(err)
	}

	defer socket.Close()

	for {
		conn, err := socket.Accept()

		if err != nil {
			log.Default().Println(err)
		}

		go n.handleConnection(conn)
	}
}

func (n *ChordServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)
	conn.Read(data)

	message := strings.Split(string(data), ",")
	opcode, err := strconv.Atoi(message[0])

	if err != nil {
		log.Default().Println(err)
	}

	responseData := ""

	switch ChordOpcode(opcode) {
	case FIND_PREDECESSOR:
		id, err := strconv.Atoi(message[1])
		if err != nil {
			log.Default().Println(err)
		}
		result := n.FindPredecessor(id)
		responseData = fmt.Sprintf("%d,%s", result.Id, result.Ip)
	case FIND_SUCCESSOR:
		id, err := strconv.Atoi(message[1])
		if err != nil {
			log.Default().Println(err)
		}
		result := n.FindSuccessor(id)
		responseData = fmt.Sprintf("%d,%s", result.Id, result.Ip)
	case GET_PREDECESSOR:
		if n.predecessor.Id != 0 {
			responseData = fmt.Sprintf("%d,%s", n.predecessor.Id, n.predecessor.Ip)
		} else {
			responseData = fmt.Sprintf("%d,%s", n.Id(), n.Reference().Ip)
		}
	case GET_SUCCESSOR:
		if n.Successor().Id != 0 {
			responseData = fmt.Sprintf("%d,%s", n.Successor().Id, n.Successor().Ip)
		} else {
			responseData = fmt.Sprintf("%d,%s", n.Id(), n.Reference().Ip)
		}
	case NOTIFY:
		ip := message[2]
		n.notify(NewChordNodeReference(ip, n.reference.Port))
	case CLOSEST_PRECEDING_FINGER:
		id, err := strconv.Atoi(message[1])
		if err != nil {
			log.Default().Println(err)
		}
		closestFinger := n.ClosestPrecedingFinger(id)
		responseData = fmt.Sprintf("%d,%s", closestFinger.Id, closestFinger.Ip)
	}

	if responseData != "" {
		conn.Write([]byte(responseData))
	}
}
