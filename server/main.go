package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	tournaments "tournament_server/tournaments"

	pb "shared/grpc"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50053, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	redisClient, err := tournaments.NewRedisClient(context.Background(), "redis:6379", "", 0)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	repo := tournaments.NewRedisRepository(redisClient)
	pb.RegisterTournamentServiceServer(s, tournaments.NewTournamentService(repo))
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// func Run() {
// 	listener, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		fmt.Println("Error creating listener:", err)
// 		os.Exit(1)
// 	}
// 	defer listener.Close()

// 	handleConnection(listener)

// 	fmt.Println("Server is listening on port 8080")

// }

// func handleConnection(listener net.Listener) {
// 	for {
// 		conn, err := listener.Accept()

// 		if err != nil {
// 			fmt.Println("Error accepting connection:", err)
// 			continue
// 		}

// 		// TODO: make it a non anonymous function
// 		go func(c net.Conn) {
// 			defer conn.Close()

// 			buffer := make([]byte, 1024)
// 			n, err := conn.Read(buffer)
// 			if err != nil {
// 				fmt.Println("Error reading from connection:", err)
// 				return
// 			}

// 			receivedString := string(buffer[:n])
// 			fmt.Println("Received string:", receivedString)
// 			// playerFactory, err := code.GetPlayerConstructor(receivedString, "NewGreedyPlayer")

// 		}(conn)
// 	}
// }
