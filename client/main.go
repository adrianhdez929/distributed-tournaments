package main

import (
	"flag"
	"log"

	"tournament_client/cli"

	pb "shared/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "10.0.11.3:50053", "the address to connect to")
)

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTournamentServiceClient(conn)

	cli.NewCliManager(client).HandleCli()
}
