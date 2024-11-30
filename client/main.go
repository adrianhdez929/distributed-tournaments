package main

import (
	"context"
	"flag"
	"log"

	pb "shared/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTournamentServiceClient(conn)

	createReq := &pb.CreateTournamentRequest{
		Name:        "New Tournament",
		Description: "A new exciting tournament",
	}

	createRes, err := client.CreateTournament(context.Background(), createReq)
	if err != nil {
		log.Fatalf("could not create tournament: %v", err)
	}
	log.Printf("Created Tournament: %v", createRes.Tournament)

	// Example: Get a tournament by ID
	getReq := &pb.GetTournamentRequest{Id: createRes.Tournament.Id}
	getRes, err := client.GetTournament(context.Background(), getReq)
	if err != nil {
		log.Fatalf("could not get tournament: %v", err)
	}
	log.Printf("Retrieved Tournament: %v", getRes.Tournament)
}
