package notifications

import (
	"context"
	"fmt"
	"log"
	"net"
	"tournament_server/chord"

	zmq "github.com/go-zeromq/zmq4"
)

const NOTIFICATION_PORT = 50055

type NotificationService struct {
	node      *chord.ChordNode
	replier   zmq.Socket
	requester zmq.Socket
	port      int
}

func NewNotificationService(node *chord.ChordNode) (*NotificationService, error) {
	// Reply server socket
	replier := zmq.NewRep(context.Background())
	requester := zmq.NewReq(context.Background())

	bindAddr := fmt.Sprintf("tcp://*:%d", NOTIFICATION_PORT)
	err := replier.Listen(bindAddr)
	if err != nil {
		replier.Close()
		log.Fatalf("Notification: failed to bind replier socket: %v", err)
		return nil, err
	}

	return &NotificationService{
		node:      node,
		replier:   replier,
		requester: requester,
		port:      NOTIFICATION_PORT,
	}, nil
}

func (s *NotificationService) Close() error {
	return s.replier.Close()
}

func (s *NotificationService) Replicate(data string, factor int) {
	successor := s.node.Client().GetSuccessor()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", successor.Ip, NOTIFICATION_PORT))
	if err != nil {
		log.Fatalf("Replicate: cannot reach successor node address")
	}

	err = s.requester.Dial(addr.String())
	if err != nil {
		log.Fatalf("Replicate: cannot connect to successor node")
	}
	defer s.requester.Close()

	const REPLICATE_OPCODE = 1
	msg := zmq.NewMsgString(fmt.Sprintf("%d;%d;%s", REPLICATE_OPCODE, factor, data))
	err = s.requester.Send(msg)
	if err != nil {
		log.Fatalf("Replicate: failed to send message to successor node")
	}

	log.Default().Printf("Replicate: replicated data to node %d", successor.Id)
}

// Notify sends a notification message to all subscribers
func (s *NotificationService) Notify(topic string, message string) error {
	// Format: topic + message (ZMQ PUB/SUB uses topic-based filtering)
	fullMessage := fmt.Sprintf("%s %s", topic, message)

	msg := zmq.NewMsgString(fullMessage)
	err := s.requester.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	log.Printf("Sent notification: %s", fullMessage)
	return nil
}
