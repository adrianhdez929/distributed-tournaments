package chord

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"tournament_server/security"
)

const REPLICATION_FACTOR = 3

const CREATE = 0
const UPDATE = 1

const MULTICAST_MASK = "224.0.0.1"
const MULTICAST_PORT = 10000

func decodeData(data []byte) string {
	decoded := strings.Replace(string(data), "\x00", "", -1)
	return decoded
}

type ChordServer struct {
	id          uint64
	successor   ChordNodeReference
	successors  []ChordNodeReference
	predecessor ChordNodeReference
	m           int
	reference   ChordNodeReference
	finger      []ChordNodeReference
	lock        *sync.Mutex
	data        map[string]string
	next        int
	channel     chan string
	// tournamentsData map[string]*pb.Tournament
}

func NewChordServer(ip string, port int, m int, serviceChannel chan string) *ChordServer {
	reference := NewChordNodeReference(ip, port)
	finger := make([]ChordNodeReference, m)

	for i := 0; i < m; i++ {
		finger[i] = reference
	}

	server := &ChordServer{
		id:          reference.Id,
		m:           m,
		successors:  make([]ChordNodeReference, m),
		successor:   reference,
		predecessor: ChordNodeReference{Id: 0, Ip: "", Port: 0},
		reference:   reference,
		finger:      finger,
		lock:        &sync.Mutex{},
		data:        make(map[string]string),
		next:        0,
		channel:     serviceChannel,
		// tournamentsData: make(map[string]*pb.Tournament),
	}

	go server.stabilize()
	go server.fixFinger()
	go server.checkPredecessor()
	go server.start()

	return server
}

func getShaRepr(data string) uint64 {
	hash := sha1.New()
	_, err := hash.Write([]byte(data))
	if err != nil {
		log.Default().Println("getShaRepr: Failed to hash data")
		return 0
	}

	hexNum := hex.EncodeToString(hash.Sum(nil))
	intNum, ok := new(big.Int).SetString(hexNum, 16)

	if !ok {
		log.Default().Println("getShaRepr: Failed to convert hex hash to int")
		return 0
	}
	return uint64(intNum.Uint64())
}

func (n *ChordServer) GetSha(data string) uint64 {
	return getShaRepr(data) % uint64(math.Pow(2, float64(n.m)))
}

func (n *ChordServer) inBetween(k uint64, start uint64, end uint64) bool {

	if start < end {
		return start < k && k <= end
	}
	return start < k || k <= end
}

func (n *ChordServer) Id() uint64 {
	return n.id
}

func (n *ChordServer) Successor() ChordNodeReference {
	return n.successor
}

func (n *ChordServer) setSuccessor(node ChordNodeReference) {
	n.lock.Lock()
	n.successor = node
	n.lock.Unlock()
}

func (n *ChordServer) Reference() ChordNodeReference {
	return n.reference
}

func (n *ChordServer) FindPredecessor(id uint64) ChordNodeReference {
	node := n.reference

	successor, err := node.Successor()

	if err != nil {
		log.Printf("FindPredecessor: cannot reach node %s", node)
		return n.reference
	}

	if successor.Id == node.Id {
		log.Printf("FindPredecessor: node %s is self successor", node)
		return node
	}

	for !n.inBetween(id, node.Id, successor.Id) && node.Id != 0 {
		var x ChordNodeReference
		var err error

		log.Default().Printf("FindPredecessor: calling ClosestPrecedingFinger from %s to %s\n", n.reference.String(), node.String())
		if node.Id == n.Id() {
			return n.ClosestPrecedingFinger(id)
		} else {
			x, err = node.ClosestPrecedingFinger(id)
		}

		if err != nil {
			log.Printf("FindPredecessor: error while calling closingPrecedingFinger to node %s", node)
			return n.reference
		}

		node = x
		successor, err = node.Successor()

		if err != nil {
			log.Printf("FindPredecessor: cannot reach node %s", node)
			return n.reference
		}

		if successor.Id == node.Id {
			log.Printf("FindPredecessor: node %s is self successor", node)
			return node
		}
	}

	return node
}

func (n *ChordServer) FindSuccessor(id uint64) ChordNodeReference {
	node := n.FindPredecessor(id)

	succ, err := node.Successor()

	if err != nil {
		return n.reference
	}

	return succ
}

func (n *ChordServer) ClosestPrecedingFinger(id uint64) ChordNodeReference {
	for i := n.m - 1; i >= 0; i-- {
		if n.finger[i].Id != 0 && n.finger[i].Id != n.reference.Id && n.inBetween(n.finger[i].Id, n.reference.Id, id) {
			return n.finger[i]
		}
	}
	return n.reference
}

func (n *ChordServer) StoreKey(key string, value string, factor int, opcode int) error {
	if factor <= 0 {
		return nil
	}

	n.lock.Lock()
	n.data[key] = value
	n.lock.Unlock()
	err := n.Successor().StoreKey(key, value, factor-1, opcode) // constant replication factor 2 this node and its successor via request

	if err != nil {
		log.Printf("storeKey: cannot replicate to successor %s", n.Successor().String())
		return err
	}

	n.channel <- fmt.Sprintf("%d;%s;%s", opcode, key, value)

	return nil
}

func (n *ChordServer) RetrieveKey(key string) (string, error) {
	kHash := n.GetSha(key)
	node := n.FindSuccessor(kHash)
	return node.RetrieveKey(key)
}

func (n *ChordServer) checkPredecessor() {
	time.Sleep(5 * time.Second)

	for {
		if n.predecessor.Id != 0 {
			err := n.predecessor.CheckPredecessor()
			if err != nil {
				log.Default().Printf("checkPredecessor: predecessor did not respond")
				n.predecessor = NewChordNodeReference("", 0)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) join(node ChordNodeReference) error {
	if node.Id != 0 {
		log.Printf("join: node %s joining node %s", n.reference.String(), node.String())
		n.predecessor = NewChordNodeReference("", 0)
		succ, err := node.FindSuccessor(n.Id())
		if err != nil {
			log.Printf("join: failed to get successor from %s", node.String())
			return err
		}
		log.Printf("join: setting node %s successor as %s", n.reference.String(), succ.String())
		n.setSuccessor(succ)
		n.Successor().Notify(n.reference)
	} else {
		n.setSuccessor(n.reference)
		n.predecessor = NewChordNodeReference("", 0)
	}

	return nil
}

func (n *ChordServer) notify(node ChordNodeReference) {
	log.Printf("notify: from %s to %s", n.reference.String(), node.String())

	if node.Id == n.Id() {
		return
	}

	if n.predecessor.Id == 0 || n.inBetween(node.Id, n.predecessor.Id, n.Id()) {
		log.Printf("notify: updating predecessor from %s to %s", n.predecessor.String(), node.String())
		n.predecessor = node
	}
}

func (n *ChordServer) stabilize() {
	time.Sleep(5 * time.Second)

	for {
		log.Printf("stabilize: predecessor %s, successor %s\n", n.predecessor.String(), n.Successor().String())
		x, err := n.Successor().Predecessor()

		if err != nil {
			log.Printf("stabilize: cannot reach successor")
			n.setSuccessor(n.reference)
		}

		if n.inBetween(x.Id, n.Id(), n.Successor().Id) && x.Id != 0 {
			log.Printf("stabilize: updating successor from %s to %s\n", n.Successor().String(), x.String())
			n.setSuccessor(x)

			// replicate keys
			go func() {
				for k, v := range n.data {
					n.StoreKey(k, v, REPLICATION_FACTOR, UPDATE)
				}
			}()
		}

		n.Successor().Notify(n.reference)
		time.Sleep(10 * time.Second)
	}
}

func (n *ChordServer) fixFinger() {
	time.Sleep(5 * time.Second)

	for {
		n.next++
		if n.next >= n.m {
			n.next = 0
		}
		node := n.FindSuccessor((n.Id() + uint64(math.Pow(2, float64(n.next)))) % uint64(math.Pow(2, float64(n.m))))
		if node.Id != 0 {
			n.finger[n.next] = node
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
		// log.Default().Printf("multicastAddress: sending address %s:%d\n", n.reference.Ip, n.reference.Port)
		conn.Close()
		time.Sleep(5 * time.Second)
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
			n.join(NewChordNodeReference(remoteIP, remotePort))
			log.Default().Printf("start: joining node %s:%d\n", parts[0], remotePort)
		}
	}

	defer mSocket.Close()
}

func (n *ChordServer) listen() {
	socket, err := security.CreateSecureSocketListener(n.reference.Port)

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

	// Read response
	var buffer bytes.Buffer
	tempBuf := make([]byte, 4096)

	for {
		n, err := conn.Read(tempBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Default().Printf("sendData: Failed to read response from %s", conn.RemoteAddr().String())
			log.Default().Println(err)
			return
		}
		buffer.Write(tempBuf[:n])

		// Check if we've reached the end of the response
		if bytes.Contains(tempBuf[:n], []byte{0}) {
			break
		}
	}

	message := strings.Split(decodeData(buffer.Bytes()), ";")
	opcode, err := strconv.Atoi(message[0])

	if err != nil {
		log.Default().Printf("handleConnection: cannot parse int from str %s\n", message[0])
		log.Default().Println(err)
	}

	responseData := ""

	switch ChordOpcode(opcode) {
	case FIND_PREDECESSOR:
		id, err := strconv.ParseUint(message[1], 10, 64)
		if err != nil {
			log.Default().Printf("handleConnection: FIND_PREDECESSOR cannot parse int from str %s\n", message[1])
			log.Default().Println(err)
		}
		result := n.FindPredecessor(id)
		responseData = fmt.Sprintf("%d;%s", result.Id, result.Ip)
	case FIND_SUCCESSOR:
		id, err := strconv.ParseUint(message[1], 10, 64)
		if err != nil {
			log.Default().Printf("handleConnection: FIND_SUCCESSOR cannot parse int from str %s\n", message[1])
			log.Default().Println(err)
		}
		result := n.FindSuccessor(id)
		responseData = fmt.Sprintf("%d;%s", result.Id, result.Ip)
	case GET_PREDECESSOR:
		if n.predecessor.Id != 0 {
			responseData = fmt.Sprintf("%d;%s", n.predecessor.Id, n.predecessor.Ip)
		} else {
			responseData = fmt.Sprintf("%d;%s", n.Id(), n.reference.Ip)
		}
	case GET_SUCCESSOR:
		if n.Successor().Id != 0 {
			responseData = fmt.Sprintf("%d;%s", n.Successor().Id, n.Successor().Ip)
		} else {
			responseData = fmt.Sprintf("%d;%s", n.Id(), n.reference.Ip)
		}
	case NOTIFY:
		ip := message[2]
		n.notify(NewChordNodeReference(ip, n.reference.Port))
	case CLOSEST_PRECEDING_FINGER:
		id, err := strconv.ParseUint(message[1], 10, 64)
		if err != nil {
			log.Default().Printf("handleConnection: CLOSEST_PRECEDING_FINGER cannot parse int from str %s\n", message[1])
			log.Default().Println(err)
		}
		closestFinger := n.ClosestPrecedingFinger(id)
		responseData = fmt.Sprintf("%d;%s", closestFinger.Id, closestFinger.Ip)
	case STORE_KEY:
		key := message[1]
		value := message[2]
		factor, err := strconv.Atoi(message[3])

		if err != nil {
			log.Printf("handleConnection: STORE_KEY cannot parse replication factor from str %s\n", message[3])
			return
		}

		opcode, err := strconv.Atoi(message[4])

		if err != nil {
			log.Printf("handleConnection: STORE_KEY cannot parse opcode from str %s\n", message[4])
			return
		}

		go n.StoreKey(key, value, factor, opcode)
	case RETRIEVE_KEY:
		key := message[1]
		value := n.data[key]
		responseData = fmt.Sprint(value)
	case CHECK_PREDECESSOR:
		responseData = "exist"
	}

	conn.Write([]byte(responseData))
}
