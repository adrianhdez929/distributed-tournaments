package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net"
	"strconv"
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

func GetSha(data string) int {
	return getShaRepr(data)
}

func checkValidIp(ip string) bool {
	ipParts := strings.Split(ip, ".")
	if len(ipParts) != 4 {
		log.Default().Printf("checkValidIp: not 4 valid parts %s\n", ip)
		return false
	}

	for _, part := range ipParts {
		partInt, err := strconv.ParseUint(string(part), 10, 64)
		if err != nil {
			log.Default().Printf("checkValidIp: cannot parse int from str %s\n in part %s\n", ip, part)
			log.Default().Println(err)
			return false
		}

		if partInt >= 255 {
			log.Default().Printf("checkValidIp: number not in range %s\n", ip)
			return false
		}
	}
	return true
}

// Stores the reference to a chord node
type ChordNodeReference struct {
	Id   int
	Ip   string
	Port int
}

func NewChordNodeReference(ip string, port int) ChordNodeReference {
	validIp := checkValidIp(ip)
	if !validIp {
		log.Default().Printf("NewChordNodeReference: Invalid IP address %s\n", ip)
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}
	return ChordNodeReference{Id: getShaRepr(ip), Ip: ip, Port: port}
}

func (n ChordNodeReference) sendData(opcode ChordOpcode, data string) []byte {
	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.Ip, n.Port))
	if err != nil {
		log.Default().Printf("sendData: Failed to connect to node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil
	}
	defer socket.Close()

	_, err = socket.Write([]byte(fmt.Sprintf("%d,%s", opcode, data)))
	if err != nil {
		log.Default().Printf("sendData: Failed to send data to node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil
	}

	if opcode == NOTIFY {
		return make([]byte, 0)
	}

	response := make([]byte, 1024)
	nBytes, err := socket.Read(response)
	if err != nil {
		log.Default().Printf("sendData: opcode %d\n", opcode)
		log.Default().Printf("sendData: Failed to read response from node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil
	}

	return response[:nBytes]
}

func (n ChordNodeReference) FindSuccessor(id int) ChordNodeReference {
	response := n.sendData(FIND_SUCCESSOR, fmt.Sprintf("%d", id))
	if response == nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	decoded := strings.Split(string(response), ",")
	log.Default().Printf("FindSuccessor: decoded %s\n", decoded[1])
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
