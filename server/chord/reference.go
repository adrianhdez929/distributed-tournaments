package chord

import (
	"fmt"
	"io"
	"log"

	"net"
	"strconv"
	"strings"
	"time"
)

type ChordOpcode int

const (
	FIND_SUCCESSOR ChordOpcode = iota + 1
	FIND_PREDECESSOR
	GET_SUCCESSOR
	GET_PREDECESSOR
	NOTIFY
	CLOSEST_PRECEDING_FINGER
	STORE_KEY
	RETRIEVE_KEY
	CHECK_PREDECESSOR
)

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
	Id   uint64
	Ip   string
	Port int
}

func NewChordNodeReference(ip string, port int) ChordNodeReference {
	if ip == "" {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}

	validIp := checkValidIp(ip)
	if !validIp {
		log.Default().Printf("NewChordNodeReference: Invalid IP address %s\n", ip)
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}
	}
	return ChordNodeReference{Id: getShaRepr(ip), Ip: ip, Port: port}
}

func (n ChordNodeReference) sendData(opcode ChordOpcode, data string) ([]byte, error) {
	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.Ip, n.Port))
	if err != nil {
		log.Default().Printf("sendData: Failed to connect to node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil, err
	}

	defer socket.Close()

	err = socket.SetDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		return nil, err
	}

	_, err = socket.Write([]byte(fmt.Sprintf("%d,%s", opcode, data)))
	if err != nil {
		log.Default().Printf("sendData: Failed to send data to node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil, err
	}

	if opcode == NOTIFY || opcode == STORE_KEY {
		return make([]byte, 0), nil
	}

	response := make([]byte, 1024)
	nBytes, err := socket.Read(response)

	if err == io.EOF {
		return make([]byte, 0), nil
	}

	if err != nil {
		log.Default().Printf("sendData: opcode %d\n", opcode)
		log.Default().Printf("sendData: Failed to read response from node %s:%d", n.Ip, n.Port)
		log.Default().Println(err)
		return nil, err
	}

	return response[:nBytes], nil
}

func (n ChordNodeReference) CheckPredecessor() error {
	_, err := n.sendData(CHECK_PREDECESSOR, "")
	return err
}

func (n ChordNodeReference) FindSuccessor(id uint64) (ChordNodeReference, error) {
	response, err := n.sendData(FIND_SUCCESSOR, fmt.Sprintf("%d", id))
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, err
	}

	decoded := strings.Split(string(response), ",")
	log.Default().Printf("FindSuccessor: decoded %s\n", decoded[1])
	return NewChordNodeReference(decoded[1], n.Port), nil
}

func (n ChordNodeReference) FindPredecessor(id int) (ChordNodeReference, error) {
	response, err := n.sendData(FIND_PREDECESSOR, fmt.Sprintf("%d", id))
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, err
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port), nil
}

func (n ChordNodeReference) Successor() (ChordNodeReference, error) {
	response, err := n.sendData(GET_SUCCESSOR, "")
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, err
	}

	parts := strings.Split(string(response), ",")
	if len(parts) != 2 {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, fmt.Errorf("failed to decode message")
	}

	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, fmt.Errorf("failed to get id from message")
	}

	return ChordNodeReference{Id: id, Ip: parts[1], Port: n.Port}, nil
}

func (n ChordNodeReference) Predecessor() (ChordNodeReference, error) {
	response, err := n.sendData(GET_PREDECESSOR, "")
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, err
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port), nil
}

func (n ChordNodeReference) Notify(node ChordNodeReference) error {
	_, err := n.sendData(NOTIFY, fmt.Sprintf("%d,%s", node.Id, node.Ip))
	return err
}

func (n ChordNodeReference) ClosestPrecedingFinger(id uint64) (ChordNodeReference, error) {
	response, err := n.sendData(CLOSEST_PRECEDING_FINGER, fmt.Sprintf("%d", id))
	if err != nil {
		return ChordNodeReference{Id: 0, Ip: "", Port: 0}, err
	}

	decoded := strings.Split(string(response), ",")
	return NewChordNodeReference(decoded[1], n.Port), nil
}

func (n ChordNodeReference) String() string {
	return fmt.Sprintf("%d:%s:%d", n.Id, n.Ip, n.Port)
}

func (n ChordNodeReference) StoreKey(key string, value string, factor int, opcode int) error {
	_, err := n.sendData(STORE_KEY, fmt.Sprintf("%s,%s,%d,%d", key, value, factor, opcode))
	return err
}

func (n ChordNodeReference) RetrieveKey(key string) (string, error) {
	response, err := n.sendData(RETRIEVE_KEY, key)

	if err != nil {
		return "", err
	}

	return decodeData(response), nil
}
