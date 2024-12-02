package tournaments

import (
	"context"
	"fmt"
	"time"
	"tournament_server/games"
	"tournament_server/models"
	"tournament_server/players"

	pb "shared/grpc"
	"shared/interfaces"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TournamentService struct {
	pb.UnimplementedTournamentServiceServer
	repo    TournamentRepository
	manager *TournamentManager
}

func NewTournamentService(repo TournamentRepository) *TournamentService {
	return &TournamentService{
		repo:    repo,
		manager: NewTournamentManager(repo),
	}
}

func (s *TournamentService) CreateTournament(ctx context.Context, req *pb.CreateTournamentRequest) (*pb.CreateTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament name is required")
	}

	tournament := s.createTournament()

	tournamentPb := &pb.Tournament{
		Id:              tournament.Id(),
		Name:            tournament.Id(),
		Description:     tournament.Id(),
		StartTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
		Status:          pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
		MaxParticipants: int32(len(tournament.Players())),
		Game:            "",
		Players:         []*pb.Player{},
	}

	if err := s.repo.Create(ctx, tournamentPb); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create tournament: %v", err)
	}

	go s.manager.AddTournament(tournament)

	return &pb.CreateTournamentResponse{
		Tournament: tournamentPb,
	}, nil
}

func (s *TournamentService) GetTournament(ctx context.Context, req *pb.GetTournamentRequest) (*pb.GetTournamentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament ID is required")
	}

	tournament, err := s.manager.GetTournament(req.Id)

	// Tournament is not in memory, fetch from database
	if err != nil {
		dbTournament, err := s.repo.Get(ctx, req.Id)

		fmt.Printf("db tournament: %v\n", dbTournament)

		fmt.Printf("requested tournament: %v\n", dbTournament)
		fmt.Printf("requested tournament status: %v\n", dbTournament.Status)
		fmt.Printf("requested tournament player wins: %v\n", dbTournament.PlayerWins)
		fmt.Printf("requested tournament final winner: %v\n", dbTournament.FinalWinner)

		if err != nil {
			return nil, status.Errorf(codes.NotFound, "tournament not found: %v", err)
		}

		return &pb.GetTournamentResponse{
			Tournament: dbTournament,
		}, nil
	}

	statistics := GetStatistics(tournament)

	fmt.Printf("requested tournament: %v\n", tournament)
	fmt.Printf("requested tournament status: %v\n", tournament.Status())
	fmt.Printf("requested tournament statistics: %v\n", statistics)
	return &pb.GetTournamentResponse{
		Tournament: &pb.Tournament{
			Id:          tournament.Id(),
			Status:      tournament.Status(),
			PlayerWins:  statistics["player_wins"].(map[string]int32),
			FinalWinner: statistics["final_winner"].(string),
		},
	}, nil
}

func (s *TournamentService) ListTournaments(ctx context.Context, req *pb.ListTournamentsRequest) (*pb.ListTournamentsResponse, error) {
	if req.PageSize <= 0 {
		req.PageSize = 50 // Default page size
	}

	tournaments, nextPageToken, err := s.repo.List(ctx, req.PageSize, req.PageToken, pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tournaments: %v", err)
	}

	return &pb.ListTournamentsResponse{
		Tournaments:   tournaments,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *TournamentService) createTournament() models.Tournament {
	// playerFactory, err := code.GetPlayerConstructor(receivedString, "NewGreedyPlayer")
	playerFactory := players.NewGreedyPlayer
	// if err != nil {
	// 	fmt.Println("Error building dynamic object:", err)
	// 	return
	// }

	// gameFactory, err := code.GetGameConstructor(receivedString, "NewTicTacToe")
	gameFactory := games.NewTicTacToe
	// if err != nil {
	// 	fmt.Println("Error building dynamic object:", err)
	// 	return
	// }

	return createTournament(playerFactory, gameFactory, 16)
}

func createTournament(playerFactory func(int) interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game, playerCount int) models.Tournament {
	players := make([]interfaces.Player, playerCount)
	// matches := make([]models.Match, playerCount/2)

	for i := 0; i < playerCount; i++ {
		players[i] = playerFactory(i + 1)
		// fmt.Printf("creating player %s\n", players[i].Id())
	}

	tournament := models.NewTournamentData(players, gameFactory)
	return tournament
}
