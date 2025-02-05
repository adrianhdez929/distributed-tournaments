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

const MULTICAST_MASK = "224.0.0.1"

func decodeData(data []byte) string {
	decoded := strings.Replace(string(data), "\x00", "", -1)
	return decoded
}

type ChordServer struct {
	id          int
	predecessor ChordNodeReference
	m           int
	reference   ChordNodeReference
	finger      []ChordNodeReference
	lock        sync.Mutex
}

func NewChordServer(ip string, port int, m int) *ChordServer {
	reference := NewChordNodeReference(ip, port)
	finger := make([]ChordNodeReference, m)

	for i := 0; i < m; i++ {
		finger[i] = reference
	}

	server := &ChordServer{
		id:          reference.Id,
		m:           m,
		predecessor: reference,
		reference:   reference,
		finger:      finger,
		lock:        sync.Mutex{},
	}

	go server.stabilize()
	go server.fixFinger()

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
	return n.finger[0]
}

func (n *ChordServer) setSuccessor(node ChordNodeReference) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.finger[0] = node
}

func (n *ChordServer) Reference() ChordNodeReference {
	return n.reference
}

func (n *ChordServer) FindPredecessor(id int) ChordNodeReference {
	node := n.reference

	if node.Id == n.Successor().Id {
		return node
	}

	for !n.inBetweenComp(id, node.Id, node.Successor().Id) && node.Id != 0 {
		log.Default().Printf("FindPredecessor: calling ClosestPrecedingFinger from %s to %s\n", n.reference.String(), node.String())
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
		if n.inRange(n.finger[i].Id, n.reference.Id, id) {
			if n.finger[i].Id != n.reference.Id {
				return n.finger[i]
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
	log.Default().Printf("join: calling FindSuccessor from %s to %s\n", n.reference.String(), node.String())
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
			log.Default().Printf("stabilize: calling GetPredecessor from %s to %s\n", n.reference.String(), n.Successor().String())
			x := n.Successor().Predecessor()

			if x.Id != n.reference.Id {
				if n.Successor().Id == n.Id() || n.inRange(x.Id, n.Id(), n.Successor().Id) {
					n.setSuccessor(x)
				}
			}

			if n.Successor().Id != 0 && n.Successor().Id != n.Id() {
				log.Default().Printf("stabilize: calling Notify from %s to %s\n", n.reference.String(), n.Successor().String())
				n.Successor().Notify(n.reference)
			}
		}
		if n.Successor().Id == 0 && n.predecessor.Id != 0 {
			n.setSuccessor(n.predecessor)
		}
		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) fixFinger() {
	time.Sleep(10 * time.Second)

	for {
		next := rand.Intn(n.m)

		n.lock.Lock()
		log.Default().Printf("fixFinger: calling FindSuccessor from %s to %s\n", n.reference.String(), n.reference.String())
		succ := n.FindSuccessor(n.reference.Id + (2^next)%(2^n.m))
		if succ.Id != 0 {
			n.finger[next] = succ
			n.lock.Unlock()

			log.Default().Printf("fixFinger:Finger table updated at index %d: %s\n", next, n.finger[next].String())
		}

		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) multicastAddress() {
	for {
		addr, err := net.ResolveUDPAddr("udp4", MULTICAST_MASK)
		if err != nil {
			log.Default().Println(err)
		}
		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			log.Default().Println(err)
		}
		conn.Write([]byte(fmt.Sprintf("%s:%d", n.reference.Ip, n.reference.Port)))
		conn.Close()
		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) start() {
	mAddr, err := net.ResolveUDPAddr("udp4", MULTICAST_MASK)

	if err != nil {
		log.Default().Println(err)
	}

	mSocket, err := net.ListenMulticastUDP("udp4", nil, mAddr)

	if err != nil {
		log.Default().Println(err)
	}

	addrReader := make([]byte, 1024)
	_, err = mSocket.Read(addrReader)

	if err != nil {
		log.Default().Println(err)
		return
	}

	// Decode the received data and remove null bytes
	receivedAddr := decodeData(addrReader)

	// Split the received address into IP and port
	parts := strings.Split(receivedAddr, ":")
	if len(parts) != 2 {
		log.Default().Printf("Invalid address format received: %s", receivedAddr)
		return
	}

	remoteIP := parts[0]
	remotePort, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Default().Printf("Invalid port number received: %s", parts[1])
		return
	}

	go n.join(NewChordNodeReference(remoteIP, remotePort))

	// Create TCP address for local listening
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.reference.Ip, n.reference.Port))

	if err != nil {
		log.Default().Println(err)
	}

	socket, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Default().Println(err)
	}

	defer mSocket.Close()
	defer socket.Close()

	go n.multicastAddress()

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

	message := strings.Split(decodeData(data), ",")
	opcode, err := strconv.Atoi(message[0])

	if err != nil {
		log.Default().Printf("handleConnection: cannot parse int from str %s\n", message[0])
		log.Default().Println(err)
	}

	responseData := ""

	switch ChordOpcode(opcode) {
	case FIND_PREDECESSOR:
		id, err := strconv.Atoi(message[1])
		if err != nil {
			log.Default().Printf("handleConnection: FIND_PREDECESSOR cannot parse int from str %s\n", message[1])
			log.Default().Println(err)
		}
		result := n.FindPredecessor(id)
		responseData = fmt.Sprintf("%d,%s", result.Id, result.Ip)
	case FIND_SUCCESSOR:
		id, err := strconv.Atoi(message[1])
		if err != nil {
			log.Default().Printf("handleConnection: FIND_SUCCESSOR cannot parse int from str %s\n", message[1])
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
			log.Default().Printf("handleConnection: CLOSEST_PRECEDING_FINGER cannot parse int from str %s\n", message[1])
			log.Default().Println(err)
		}
		closestFinger := n.ClosestPrecedingFinger(id)
		responseData = fmt.Sprintf("%d,%s", closestFinger.Id, closestFinger.Ip)
	}

	conn.Write([]byte(responseData))
}
