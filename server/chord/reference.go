package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net"
	"strings"
)

type ChordOpcode int

const (
	FIND_SUCCESSOR ChordOpcode = iota + 1
	FIND_PREDECESSOR
	GET_SUCCESSOR
	GET_PREDECESSOR
	NOTIFY
	CLOSEST_PRECEDING_FINGER
)

func getShaRepr(data string) int {
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
	return int(intNum.Int64())
}

// Stores the reference to a chord node
type ChordNodeReference struct {
	Id   int
	Ip   string
	Port int
}

func NewChordNodeReference(ip string, port int) ChordNodeReference {
	return ChordNodeReference{Id: getShaRepr(ip), Ip: ip, Port: port}
}

func (n ChordNodeReference) sendData(opcode ChordOpcode, data string) []byte {
	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.Ip, n.Port))
	if err != nil {
		log.Default().Printf("sendData: Failed to connect to node %s:%d", n.Ip, n.Port)
		return nil
	}
	defer socket.Close()

	_, err = socket.Write([]byte(fmt.Sprintf("%d,%s", opcode, data)))
	if err != nil {
		log.Default().Printf("sendData: Failed to send data to node %s:%d", n.Ip, n.Port)
		return nil
	}

	response := make([]byte, 1024)
	_, err = socket.Read(response)
	if err != nil {
		log.Default().Printf("sendData: Failed to read response from node %s:%d", n.Ip, n.Port)
		return nil
	}

	return response
}

func (n ChordNodeReference) FindSuccessor(id int) ChordNodeReference {
	response := n.sendData(FIND_SUCCESSOR, fmt.Sprintf("%d", id))
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port)
}

func (n ChordNodeReference) FindPredecessor(id int) ChordNodeReference {
	response := n.sendData(FIND_PREDECESSOR, fmt.Sprintf("%d", id))
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port)
}

func (n *ChordNodeReference) Successor() ChordNodeReference {
	response := n.sendData(GET_SUCCESSOR, "")
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port)
}

func (n ChordNodeReference) Predecessor() ChordNodeReference {
	response := n.sendData(GET_PREDECESSOR, "")
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port)
}

func (n ChordNodeReference) Notify(node ChordNodeReference) {
	n.sendData(NOTIFY, fmt.Sprintf("%d,%s", node.Id, node.Ip))
}

func (n ChordNodeReference) ClosestPrecedingFinger(id int) ChordNodeReference {
	response := n.sendData(CLOSEST_PRECEDING_FINGER, fmt.Sprintf("%d", id))
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port)
}

func (n ChordNodeReference) String() string {
	return fmt.Sprintf("%d:%s:%d", n.Id, n.Ip, n.Port)
}
