package tournaments

import (
	"math/rand"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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

const REPLICATION_FACTOR = 3

const CREATE = 0
const UPDATE = 1

type TournamentService struct {
	pb.UnimplementedTournamentServiceServer
	repo    TournamentRepository
	manager *TournamentManager
	node    *chord.ChordServer
	channel chan string
}

func NewTournamentService(repo TournamentRepository, manager *TournamentManager, node *chord.ChordServer, channel chan string) *TournamentService {
	service := &TournamentService{
		repo:    repo,
		manager: manager,
		node:    node,
		channel: channel,
	}
	go service.replicateData()
	go service.handleNotifications()

	return service
}

func (s *TournamentService) parseNotification(notification string) {
	parts := strings.Split(notification, ",")
	opcode, err := strconv.Atoi(parts[0])

	if err != nil {
		log.Printf("tournamentService:parseNotification: cannot parse opcode from str %s\n", parts[0])
	}

	switch opcode {
	case 0: // REPLICATE_CREATE
		key := parts[1]
		owner := parts[2]

		tournamentId := strings.Split(key, ":")[1]

		conn, err := grpc.NewClient(
			fmt.Sprintf("%s:%d", owner, 50053),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			log.Printf("tournamentService:parseNotification: failed to connect to key owner")
			return
		}

		client := pb.NewTournamentServiceClient(conn)
		_, err = client.CreateTournament(context.Background(), &pb.CreateTournamentRequest{Name: tournamentId})

		if err != nil {
			log.Printf("tournamentService:parseNotification: could not update tournament in %s", owner)
			return
		}
	case 1: // REPLICATE_UPDATE
		key := parts[1]
		owner := parts[2]

		tournamentId := strings.Split(key, ":")[1]
		tournament, _ := s.repo.Get(context.Background(), tournamentId)

		conn, err := grpc.NewClient(
			fmt.Sprintf("%s:%d", owner, 50053),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			log.Printf("tournamentService:parseNotification: failed to connect to key owner")
			return
		}

		client := pb.NewTournamentServiceClient(conn)
		_, err = client.UpdateTournament(context.Background(), &pb.UpdateTournamentRequest{
			Tournament: tournament,
		})
		if err != nil {
			log.Printf("tournamentService:parseNotification: could not update tournament in %s", owner)
			return
		}
	}
}

func (s *TournamentService) handleNotifications() {
	for {
		// this sends ip addresses to replicate data to
		notification := <-s.channel
		s.parseNotification(notification)

	}
}

func (s *TournamentService) replicateData() {
	for {
		log.Printf("tournamentService:replicateData: replicating tournament data")
		tList, err := s.repo.List(context.Background())
		if err != nil {
			return
		}

		for _, t := range tList {
			owner := s.getTournamentOwner(t.Name)
			log.Printf("tournamentService:replicateData: replicating tournament %s to node %d", t.Name, owner.Id)
			conn, err := grpc.NewClient(
				fmt.Sprintf("%s:%d", owner.Ip, 50053),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			if err != nil {
				log.Printf("failed to connect to connect to key owner")
				continue
			}

			client := pb.NewTournamentServiceClient(conn)
			client.UpdateTournament(context.Background(), &pb.UpdateTournamentRequest{Tournament: t})
		}

		time.Sleep(20 * time.Second)
	}
}

func (s *TournamentService) getTournamentKey(name string) string {
	return fmt.Sprintf("tournament:%s", name)
}

func (s *TournamentService) getTournamentOwner(name string) chord.ChordNodeReference {
	tHash := chord.GetSha(s.getTournamentKey(name))
	owner := s.node.FindSuccessor(tHash)
	return owner
}

func (s *TournamentService) CreateTournament(ctx context.Context, req *pb.CreateTournamentRequest) (*pb.CreateTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament name is required")
	}

	log.Default().Printf("request to create tournament in service")

	owner := s.getTournamentOwner(req.Name)

	log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

	if owner.Id == 0 {
		return nil, status.Error(codes.NotFound, "tournament not found")
	}

	value, err := owner.RetrieveKey(s.getTournamentKey(req.Name))

	if err != nil {
		log.Default().Fatalf("tournamentService:createTournament: there was an error getting the tournament key")
	}

	if value == "" {
		// this server is the owner of the resource
		if owner.Ip == s.node.Reference().Ip {
			log.Default().Printf("owner is this server, creating tournament")
			tournament := s.createTournament(req.Name)

			tournamentPb := &pb.Tournament{
				Id:              tournament.Id(),
				Name:            req.Name,
				Description:     req.Description,
				StartTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
				Status:          pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
				MaxParticipants: int32(len(tournament.Players())),
				Game:            tournament.Game(),
				Players:         DumpTournamentPlayers(tournament.Players()),
				Matches:         DumpTournamentMatches(tournament.Matches()),
			}

			if err := s.repo.Create(ctx, tournamentPb); err != nil {
				log.Fatalf("failed to create tournament: %s", err)
				return nil, status.Errorf(codes.Internal, "failed to create tournament: %v", err)
			}

			go s.manager.AddTournament(tournament)

			err := owner.StoreKey(s.getTournamentKey(req.Name), "existe", REPLICATION_FACTOR, CREATE)
			if err != nil {
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

func (s *TournamentService) ReplicateCreateRequest(ctx context.Context, req *pb.ReplicateCreateRequest) (*pb.CreateTournamentResponse, error) {
	if req.Tournament.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament name is required")
	}

	_, err := s.manager.GetTournament(req.Tournament.Name)

	if err != nil {
		// tournament is not already running
		tournament := s.createTournament(req.Tournament.Name)
		go s.manager.AddTournament(tournament)

		tournamentPb := &pb.Tournament{
			Id:              tournament.Id(),
			Name:            req.Tournament.Name,
			Description:     req.Tournament.Description,
			StartTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
			Status:          pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED,
			MaxParticipants: int32(len(tournament.Players())),
			Game:            tournament.Game(),
			Players:         DumpTournamentPlayers(tournament.Players()),
			Matches:         DumpTournamentMatches(tournament.Matches()),
		}

		if err := s.repo.Create(ctx, tournamentPb); err != nil {
			log.Fatalf("failed to create tournament: %s", err)
			return nil, status.Errorf(codes.Internal, "failed to create tournament: %v", err)
		}

		return &pb.CreateTournamentResponse{
			Tournament: tournamentPb,
		}, nil
	} else {
		return nil, status.Error(codes.AlreadyExists, "tournament already exists")
	}
}

func (s *TournamentService) GetTournament(ctx context.Context, req *pb.GetTournamentRequest) (*pb.GetTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament Name is required")
	}

	log.Default().Printf("request to get tournament in service")

	log.Default().Printf("tournamentService:get: checking if the server has a copy of %s", req.Name)

	tournament, _ := s.repo.Get(ctx, req.Name)

	if tournament != nil {
		log.Default().Printf("tournamentService:get: the server has a copy of %s", req.Name)
		return &pb.GetTournamentResponse{
			Tournament: tournament,
		}, nil
	}

	owner := s.getTournamentOwner(req.Name)
	log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

	if owner.Id == 0 {
		return nil, status.Error(codes.NotFound, "tournament not found")
	}

	value, err := owner.RetrieveKey(s.getTournamentKey(req.Name))

	if err != nil {
		return nil, status.Error(codes.Internal, "could not retrieve the key from owner")
	}

	if value == "" {
		return nil, status.Error(codes.NotFound, "tournament not found")
	} else {
		// this server is the owner of the resource
		if owner.Ip == s.node.Reference().Ip {

			tournament, err := s.repo.Get(context.Background(), req.Name)

			if err != nil {
				return nil, status.Error(codes.Internal, "could not get tournament from repo")
			}

			return &pb.GetTournamentResponse{
				Tournament: tournament,
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

func (s *TournamentService) UpdateTournament(ctx context.Context, req *pb.UpdateTournamentRequest) (*pb.UpdateTournamentResponse, error) {
	err := s.repo.Update(ctx, req.Tournament)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update tournament")
	}
	s.node.StoreKey(s.getTournamentKey(req.Tournament.Name), "existe", REPLICATION_FACTOR, UPDATE)

	return &pb.UpdateTournamentResponse{Success: true}, nil
}

func (s *TournamentService) ListTournaments(ctx context.Context, req *pb.ListTournamentsRequest) (*pb.ListTournamentsResponse, error) {
	if req.PageSize <= 0 {
		req.PageSize = 50 // Default page size
	}

	tournaments, err := s.repo.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tournaments: %v", err)
	}

	return &pb.ListTournamentsResponse{
		Tournaments: tournaments,
	}, nil
}

func (s *TournamentService) createTournament(id string) models.Tournament {
	// playerFactory, err := code.GetPlayerConstructor(receivedString, "NewGreedyPlayer")
	playerFactory := players.NewRandomPlayer
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
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})
	

	tournament := models.NewTournamentData(id, players, gameFactory)
	return tournament
}
