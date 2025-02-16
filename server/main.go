package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "shared/grpc"
	chord "tournament_server/chord"
	persistency "tournament_server/persistency"
	tournaments "tournament_server/tournaments"

	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 50053, "The server port")
	chord_port = flag.Int("chord_port", 50054, "The chord port")
	ip         = flag.String("ip", "10.0.11.3", "The ip address")
)

// func mainMonolithic() {
// 	flag.Parse()
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	redisClient, err := persistency.NewRedisClient(context.Background(), "redis:6379", "", 0)
// 	if err != nil {
// 		log.Fatalf("failed to connect to Redis: %v", err)
// 	}
// 	repo := tournaments.NewRedisRepository(redisClient)
// 	pb.RegisterTournamentServiceServer(s, tournaments.NewTournamentService(repo))
// 	log.Printf("server listening at %v", lis.Addr())

// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }

func main() {
	flag.Parse()

	node := chord.NewChordNode(*ip, *chord_port)

	// Tournament client handler
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	redisClient, err := persistency.NewRedisClient(context.Background(), "redis:6379", "", 0)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	repo := tournaments.NewRedisRepository(redisClient)
	manager := tournaments.NewTournamentManager(repo)
	pb.RegisterTournamentServiceServer(s, tournaments.NewTournamentService(repo, manager, node))
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// fmt.Println(chord.NewChordNodeReference("0", 50054).String())
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
