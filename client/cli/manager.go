package cli

import (
	"flag"
	"fmt"
	"os"
	pb "shared/grpc"
)

type CliManager struct {
	client pb.TournamentServiceClient
}

func NewCliManager(client pb.TournamentServiceClient) *CliManager {
	return &CliManager{client: client}
}

func (m *CliManager) HandleCli() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("expected 'create', 'get', or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		CreateTournament(m.client)
	case "get":
		GetTournament(m.client)
	case "list":
		ListTournaments(m.client)
	default:
		fmt.Println("expected 'create', 'get', or 'list' subcommands")
		os.Exit(1)
	}
}
