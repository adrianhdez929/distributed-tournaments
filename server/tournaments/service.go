package tournaments

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	"google.golang.org/protobuf/encoding/protojson"
)

func TournamentToJson(tournament *pb.Tournament) (string, error) {
	json, err := protojson.Marshal(tournament)
	if err != nil {
		log.Printf("failed to marshal tournament")
		return "nil", fmt.Errorf("failed to marshal tournament")
	}
	return string(json), nil
}

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	parts := strings.Split(notification, ";")
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
		_, err = client.CreateTournament(ctx, &pb.CreateTournamentRequest{Name: tournamentId})

		if err != nil {
			log.Printf("tournamentService:parseNotification: could not create tournament in %s", owner)
			log.Printf("%s", err)
			return
		}
	case 1: // REPLICATE_UPDATE
		key := parts[1]
		value := parts[2]

		log.Printf("tornamentService:parseNotification: received update for tournament %s", key)

		tournament := models.TournamentData{}.FromJson(value)
		local, _ := s.repo.Get(context.Background(), tournament.Id())

		if local != nil {
			if local.Status < tournament.Status() || local.Status == pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS && local.Status == tournament.Status() {
				s.updateTournamentFromJson(value)
			}
			return
		}

		s.updateTournamentFromJson(value)
	}
}

func (s *TournamentService) handleNotifications() {
	for {
		// this sends ip addresses to replicate data to
		for notification := range s.channel {
			s.parseNotification(notification)
		}
	}
}

func (s *TournamentService) replicateData() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		log.Printf("tournamentService:replicateData: replicating tournament data")
		tList, err := s.repo.List(ctx)
		if err != nil {
			cancel()
			time.Sleep(20 * time.Second)
			continue
		}

		for _, t := range tList {
			owner := GetTournamentOwner(s.node, t.Name)
			log.Printf("tournamentService:replicateData: replicating tournament %s to node %d", t.Name, owner.Id)

			json, err := TournamentToJson(t)
			if err != nil {
				log.Println(err)
				continue
			}

			err = s.node.StoreKey(GetTournamentKey(t.Name), string(json), chord.REPLICATION_FACTOR, chord.UPDATE)

			if err != nil {
				log.Printf("failed to connect to connect to key owner")
				log.Printf("resuming tournament in this node")

				json, err := TournamentToJson(t)
				if err != nil {
					log.Printf("failed to marshal tournament")
					continue
				}
				if t.Status == pb.TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS {
					go s.manager.ResumeTournament(json)
				}
				continue
			}
		}

		cancel()
		time.Sleep(20 * time.Second)
	}
}

func (s *TournamentService) CreateTournament(ctx context.Context, req *pb.CreateTournamentRequest) (*pb.CreateTournamentResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tournament name is required")
	}

	log.Default().Printf("request to create tournament in service")

	owner := GetTournamentOwner(s.node, req.Name)

	log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

	if owner.Id == 0 {
		return nil, status.Error(codes.NotFound, "tournament not found")
	}

	value, err := owner.RetrieveKey(GetTournamentKey(req.Name))

	if err != nil {
		log.Default().Fatalf("tournamentService:createTournament: there was an error getting the tournament key")
	}

	if value == "" {
		// this server is the owner of the resource
		if owner.Ip == s.node.Reference().Ip {
			log.Default().Printf("owner is this server, creating tournament")
			tournament := s.createTournament(req.Name)

			tournamentPb := DumpTournament(tournament)

			if err := s.repo.Create(ctx, tournamentPb); err != nil {
				log.Fatalf("failed to create tournament: %s", err)
				return nil, status.Errorf(codes.Internal, "failed to create tournament: %v", err)
			}

			go s.manager.StartTournament(tournament)

			json, err := TournamentToJson(tournamentPb)
			if err != nil {
				return nil, status.Error(codes.Internal, "error while serializing tournament data")
			}

			err = owner.StoreKey(GetTournamentKey(req.Name), json, chord.REPLICATION_FACTOR, chord.UPDATE)
			if err != nil {
				return nil, status.Error(codes.Internal, "error while storing key")
			}
			log.Default().Printf("tournamentService:createTournament: the key %s value is %s", GetTournamentKey(req.Name), value)

			s.updateTournamentFromJson(json)

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
		log.Default().Printf("tournamentService:createTournament: there is an existent key %s with value %s", GetTournamentKey(req.Name), value)
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
		go s.manager.StartTournament(tournament)

		tournamentPb := DumpTournament(tournament)

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

	t, err := s.repo.Get(context.Background(), req.Name)

	// hay copia local
	if err == nil && t != nil {
		return &pb.GetTournamentResponse{
			Tournament: t,
		}, nil
	} else {
		owner := GetTournamentOwner(s.node, req.Name)
		log.Default().Printf("tournamentService:get: owner of tournament %s is %s", req.Name, owner.Ip)

		if owner.Id == 0 {
			return nil, status.Error(codes.NotFound, "tournament not found")
		}

		value, err := owner.RetrieveKey(GetTournamentKey(req.Name))

		if err != nil {
			return nil, status.Error(codes.Internal, "could not retrieve the key from owner")
		}

		if value == "" {
			return nil, status.Error(codes.NotFound, "tournament not found")
		} else {
			// this server is the owner of the resource
			if owner.Ip == s.node.Reference().Ip {
				_, err := s.manager.GetStatus(req.Name)

				if err != nil {
					tournament, err := s.repo.Get(context.Background(), req.Name)

					if err != nil {
						return nil, status.Error(codes.Internal, "could not get tournament from repo")
					}

					json, err := TournamentToJson(tournament)

					if err != nil {
						return nil, status.Error(codes.Internal, "could not deserialize tournament data")
					}

					go s.manager.ResumeTournament(json)
				}

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
}

func (s *TournamentService) UpdateTournament(ctx context.Context, req *pb.UpdateTournamentRequest) (*pb.UpdateTournamentResponse, error) {
	err := s.repo.Update(ctx, req.Tournament)
	if err != nil {
		log.Printf("tournamentService:UpdateTournament: repo error: %s", err)
		return nil, status.Error(codes.Internal, "could not update tournament")
	}

	json, err := TournamentToJson(req.Tournament)
	if err != nil {
		return nil, status.Error(codes.Internal, "error while serializing tournament data")
	}

	go s.node.StoreKey(GetTournamentKey(req.Tournament.Name), json, chord.REPLICATION_FACTOR, chord.UPDATE)
	s.updateTournamentFromJson(json)

	return &pb.UpdateTournamentResponse{Success: true}, nil
}

func (s *TournamentService) updateTournamentFromJson(data string) {
	tournament := models.TournamentData{}.FromJson(data)
	if tournament == nil {
		log.Printf("failed to load tournament from json")
		return
	}

	tournamentPb := DumpTournament(tournament)

	t, _ := s.repo.Get(context.Background(), tournament.Id())

	if t != nil {
		if err := s.repo.Update(context.Background(), tournamentPb); err != nil {
			log.Printf("failed to update tournament: %s", err)
			return
		}
	} else {
		s.repo.Create(context.Background(), tournamentPb)
	}

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

func createTournament(id string, playerFactory func(string) interfaces.Player, gameFactory func([]interfaces.Player) interfaces.Game, playerCount int) models.Tournament {
	players := make([]interfaces.Player, playerCount)
	// matches := make([]models.Match, playerCount/2)

	for i := 0; i < playerCount; i++ {
		players[i] = playerFactory(strconv.Itoa(i + 1))
		// fmt.Printf("creating player %s\n", players[i].Id())
	}
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	tournament := models.NewTournamentData(id, players, gameFactory)
	return tournament
}
