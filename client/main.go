package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"tournament_client/cli"

	pb "shared/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MULTICAST_MASK = "224.0.0.1"
const MULTICAST_PORT = 10000

func decodeData(data []byte) string {
	decoded := strings.Replace(string(data), "\x00", "", -1)
	return decoded
}

func main() {
	mAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", MULTICAST_MASK, MULTICAST_PORT))

	if err != nil {
		log.Default().Println(err)
	}

	log.Default().Println("Listening to multicast address")
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
	log.Default().Printf("start: received address %s\n", receivedAddr)

	parts := strings.Split(receivedAddr, ":")
	if len(parts) != 2 {
		log.Default().Printf("Invalid address format received: %s", receivedAddr)
		return
	}

	// Split the received address into IP and port
	remoteIP := parts[0]
	// remoteIP := "10.0.11.2"
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", remoteIP, 50053), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTournamentServiceClient(conn)

	cli.NewCliManager(client).HandleCli()
}
