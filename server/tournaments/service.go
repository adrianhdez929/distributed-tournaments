package tournaments

import (
	"context"
	"fmt"
	"log"
	"time"
	"tournament_server/chord"
	"tournament_server/games"
	"tournament_server/models"
	"tournament_server/players"

	pb "shared/grpc"
	"shared/interfaces"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type TournamentService struct {
	pb.UnimplementedTournamentServiceServer
	repo    TournamentRepository
	manager *TournamentManager
	node    *chord.ChordNode
}

func NewTournamentService(repo TournamentRepository, manager *TournamentManager, node *chord.ChordNode) *TournamentService {
	return &TournamentService{
		repo:    repo,
		manager: manager,
		node:    node,
	}
}

func (s *TournamentService) getTournamentKey(name string) string {
	return fmt.Sprintf("tournament:%s", name)
}

func (s *TournamentService) getTournamentOwner(name string) chord.ChordNodeReference {
	tHash := chord.GetSha(s.getTournamentKey(name))
	owner := s.node.Server().FindSuccessor(tHash)
	return owner
}

func (s *TournamentService) CreateTournament(ctx context.Context, req *pb.CreateTournamentRequest) (*pb.CreateTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament name is required")
	}

	log.Default().Printf("request to create tournament in service")

	owner := s.getTournamentOwner(req.Name)
	log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

	value, err := owner.RetrieveKey(s.getTournamentKey(req.Name))

	if err != nil {
		log.Default().Fatalf("tournamentService:createTournament: there was an error getting the tournament key")
	}

	if value == "" {
		// this server is the owner of the resource
		if owner.Ip == s.node.Server().Reference().Ip {
			log.Default().Printf("owner is this server, creating tournament")
			tournament := s.createTournament(req.Name)

			tournamentPb := &pb.Tournament{
				Id:              tournament.Id(),
				Name:            req.Name,
				Description:     req.Description,
				StartTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
				Status:          pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
				MaxParticipants: int32(len(tournament.Players())),
				Game:            "",
				Players:         []*pb.Player{},
			}

			if err := s.repo.Create(ctx, tournamentPb); err != nil {
				log.Fatalf("failed to create tournament: %s", err)
				return nil, status.Errorf(codes.Internal, "failed to create tournament: %v", err)
			}

			go s.manager.AddTournament(tournament)

			result := owner.StoreKey(s.getTournamentKey(req.Name), "existe, no se que guardar aqui ahora mimo")
			if result == nil {
				return nil, status.Error(codes.Internal, "error while storing key")
			}
			log.Default().Printf("tournamentService:createTournament: the key %s value is %s", s.getTournamentKey(req.Name), value)

			return &pb.CreateTournamentResponse{
				Tournament: tournamentPb,
			}, nil
			// the owner is another server
		} else {
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", owner.Ip, 50053),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			if err != nil {
				return nil, status.Error(codes.Internal, "failed to connect to connect to key owner")
			}

			client := pb.NewTournamentServiceClient(conn)
			res, err := client.CreateTournament(ctx, req)

			return res, err
		}

	} else {
		log.Default().Printf("tournamentService:createTournament: there is an existent key %s with value %s", s.getTournamentKey(req.Name), value)
		return nil, status.Error(codes.AlreadyExists, "tournament already exists")
	}
}

func (s *TournamentService) GetTournament(ctx context.Context, req *pb.GetTournamentRequest) (*pb.GetTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament Name is required")
	}

	log.Default().Printf("request to get tournament in service")

	owner := s.getTournamentOwner(req.Name)
	log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

	value, err := owner.RetrieveKey(s.getTournamentKey(req.Name))

	if err != nil {
		return nil, status.Error(codes.Internal, "could not retrieve the key from owner")
	}

	if value == "" {
		return nil, status.Error(codes.NotFound, "tournament not found")
	} else {
		// this server is the owner of the resource
		if owner.Ip == s.node.Server().Reference().Ip {
			tournament, err := s.manager.GetTournament(req.Name)
			log.Default().Printf("tournamentService:get: owner of %s is this server %s", req.Name, owner.Ip)
			log.Default().Printf("tournamentService:get: retrieving tournament %s with key %s", req.Name, s.getTournamentKey(req.Name))

			if err != nil {
				return nil, status.Error(codes.Internal, "could not retrieve tournament from manager")
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
					FinalWinner: statistics["winner"].(interfaces.Player).Id(),
				},
			}, nil
			// the owner is another server in the network
		} else {
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", owner.Ip, 50053),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			if err != nil {
				return nil, status.Error(codes.Internal, "failed to connect to connect to key owner")
			}

			client := pb.NewTournamentServiceClient(conn)
			res, err := client.GetTournament(ctx, &pb.GetTournamentRequest{Name: req.Name})

			return res, err
		}
	}
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

func (s *TournamentService) createTournament(id string) models.Tournament {
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

	return createTournament(id, playerFactory, gameFactory, 16)
}

func createTournament(id string, playerFactory func(int) interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game, playerCount int) models.Tournament {
	players := make([]interfaces.Player, playerCount)
	// matches := make([]models.Match, playerCount/2)

	for i := 0; i < playerCount; i++ {
		players[i] = playerFactory(i + 1)
		// fmt.Printf("creating player %s\n", players[i].Id())
	}

	tournament := models.NewTournamentData(id, players, gameFactory)
	return tournament
}
