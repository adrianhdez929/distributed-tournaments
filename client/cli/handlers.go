package cli

import (
	"context"
	"flag"
	"log"
	"os"
	pb "shared/grpc"
	"time"
)

func CreateTournament(client pb.TournamentServiceClient) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	name := createCmd.String("name", "", "Name of the tournament")
	description := createCmd.String("description", "", "Description of the tournament")

	createCmd.Parse(os.Args[2:])

	if *name == "" {
		log.Fatal("Name is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createReq := &pb.CreateTournamentRequest{
		Name:        *name,
		Description: *description,
	}

	createRes, err := client.CreateTournament(ctx, createReq)
	if err != nil {
		log.Fatalf("could not create tournament: %v", err)
	}
	log.Printf("%s", createRes.Tournament.Id)
}

func GetTournament(client pb.TournamentServiceClient) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	name := getCmd.String("name", "", "Name of the tournament")

	getCmd.Parse(os.Args[2:])

	if *name == "" {
		log.Fatal("Name is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getReq := &pb.GetTournamentRequest{Name: *name}
	getRes, err := client.GetTournament(ctx, getReq)
	if err != nil {
		log.Fatalf("could not get tournament: %v", err)
	}
	log.Printf("%v", getRes.Tournament)
}

func ListTournaments(client pb.TournamentServiceClient) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	pageSize := listCmd.Int("page_size", 10, "Number of tournaments to list")
	pageToken := listCmd.String("page_token", "", "Page token for pagination")

	listCmd.Parse(os.Args[2:])

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	listReq := &pb.ListTournamentsRequest{
		PageSize:  int32(*pageSize),
		PageToken: *pageToken,
	}

	listRes, err := client.ListTournaments(ctx, listReq)
	if err != nil {
		log.Fatalf("could not list tournaments: %v", err)
	}
	for _, tournament := range listRes.Tournaments {
		log.Printf("Tournament: %v", tournament)
	}
	if listRes.NextPageToken != "" {
		log.Printf("Next Page Token: %s", listRes.NextPageToken)
	}
}
