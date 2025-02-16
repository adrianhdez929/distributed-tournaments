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
const MULTICAST_PORT = 10000

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
	lock        *sync.Mutex
	data        map[string]string
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
		lock:        &sync.Mutex{},
		data:        make(map[string]string),
	}

	go server.stabilize()
	go server.fixFinger()

	go server.start()

	return server
}

func (n *ChordServer) inBetween(k int, start int, end int) bool {
	if start < end {
		return start < k && k <= end
	} else {
		return start < k || k <= end
	}
}

func (n *ChordServer) Id() int {
	return n.id
}

func (n *ChordServer) Successor() ChordNodeReference {
	return n.finger[0]
}

func (n *ChordServer) setSuccessor(node ChordNodeReference) {
	n.lock.Lock()
	n.finger[0] = node
	n.lock.Unlock()
}

func (n *ChordServer) Reference() ChordNodeReference {
	return n.reference
}

func (n *ChordServer) FindPredecessor(id int) ChordNodeReference {
	node := n.reference

	if node.Id == n.Successor().Id {
		return node
	}

	for !n.inBetween(id, node.Id, node.Successor().Id) && node.Id != 0 {
		log.Default().Printf("FindPredecessor: calling ClosestPrecedingFinger from %s to %s\n", n.reference.String(), node.String())
		node = node.ClosestPrecedingFinger(id)

		// if node.Id == n.reference.Id {
		// 	break
		// }
	}

	return node
}

func (n *ChordServer) FindSuccessor(id int) ChordNodeReference {
	node := n.FindPredecessor(id)

	return node.Successor()
}

func (n *ChordServer) ClosestPrecedingFinger(id int) ChordNodeReference {
	for i := n.m - 1; i >= 0; i-- {
		if n.finger[i].Id != n.reference.Id && n.inBetween(n.finger[i].Id, n.reference.Id, id) {
			return n.finger[i]
		}
	}
	return n.reference
}

func (n *ChordServer) StoreKey(key string, value string) error {
	kHash := getShaRepr(key)
	node := n.FindSuccessor(kHash)
	response := node.StoreKey(key, value)

	if response == nil {
		return fmt.Errorf("could not store key on node %s", node.Ip)
	}

	n.data[key] = value
	n.Successor().StoreKey(key, value) // constant replication factor 2 this node and its successor via request
	return nil
}

func (n *ChordServer) RetrieveKey(key string) (string, error) {
	kHash := getShaRepr(key)
	node := n.FindSuccessor(kHash)
	return node.RetrieveKey(key)
}

func (n *ChordServer) join(node ChordNodeReference) ChordNodeReference {
	time.Sleep(10 * time.Second)
	n.predecessor = n.reference
	log.Default().Printf("join: calling FindSuccessor from %s to %s\n", n.reference.String(), node.String())
	successor := node.FindSuccessor(n.Id())
	if successor.Id != 0 {
		n.setSuccessor(successor)
	}
	return successor
}

func (n *ChordServer) notify(node ChordNodeReference) {
	if node.Id == n.reference.Id {
		return
	}

	if n.predecessor.Id == n.reference.Id || n.inBetween(node.Id, n.predecessor.Id, n.reference.Id) {
		n.predecessor = node
	}
}

func (n *ChordServer) stabilize() {
	for {
		if n.Successor().Id != 0 && n.Successor().Id != n.Reference().Id {
			log.Default().Printf("stabilize: calling GetPredecessor from %s to %s\n", n.reference.String(), n.Successor().String())
			x := n.Successor().Predecessor()

			if x.Id != n.reference.Id {
				if n.inBetween(x.Id, n.reference.Id, n.Successor().Id) {
					n.setSuccessor(x)
				}
				n.Successor().Notify(n.reference)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (n *ChordServer) fixFinger() {
	time.Sleep(5 * time.Second)

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
		addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", MULTICAST_MASK, MULTICAST_PORT))
		if err != nil {
			log.Default().Println(err)
		}
		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			log.Default().Println(err)
		}
		conn.Write([]byte(fmt.Sprintf("%s:%d", n.reference.Ip, n.reference.Port)))
		log.Default().Printf("multicastAddress: sending address %s:%d\n", n.reference.Ip, n.reference.Port)
		conn.Close()
		time.Sleep(30 * time.Second)
	}
}

func (n *ChordServer) start() {
	go n.multicastAddress()
	go n.listen()

	mAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", MULTICAST_MASK, MULTICAST_PORT))

	if err != nil {
		log.Default().Println(err)
	}

	mSocket, err := net.ListenMulticastUDP("udp4", nil, mAddr)

	if err != nil {
		log.Default().Println(err)
	}

	addressFound := false
	remoteIP := ""
	remotePort := 0

	for !addressFound {
		addrReader := make([]byte, 1024)
		_, err = mSocket.Read(addrReader)

		if err != nil {
			log.Default().Println(err)
			return
		}

		// Decode the received data and remove null bytes
		receivedAddr := decodeData(addrReader)
		log.Default().Printf("start: received address %s\n", receivedAddr)

		parts := strings.Split(receivedAddr, ":")
		if len(parts) != 2 {
			log.Default().Printf("Invalid address format received: %s", receivedAddr)
			return
		}

		// Split the received address into IP and port
		remoteIP = parts[0]
		remotePort, err = strconv.Atoi(parts[1])

		if err != nil {
			log.Default().Printf("Invalid port number received: %s", parts[1])
			return
		}

		if parts[0] != n.reference.Ip {
			addressFound = true
			newReference := n.join(NewChordNodeReference(remoteIP, remotePort))
			log.Default().Printf("start: joining node %s:%d\n", parts[0], remotePort)

			if newReference.Id == 0 {
				addressFound = false
				log.Default().Fatalf("start: failed to join node %s:%d\n", parts[0], remotePort)
			}
		}
	}

	defer mSocket.Close()
}

func (n *ChordServer) listen() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", n.reference.Ip, n.reference.Port))

	if err != nil {
		log.Default().Println(err)
	}

	socket, err := net.ListenTCP("tcp", addr)

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
	case STORE_KEY:
		key := message[1]
		value := message[2]
		n.data[key] = value
		// hay que replicar la ejecucion
	case RETRIEVE_KEY:
		key := message[1]
		responseData = fmt.Sprint(n.data[key])
	}

	conn.Write([]byte(responseData))
}
