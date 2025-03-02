package tournaments

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "shared/grpc"

	"github.com/redis/go-redis/v9"
)

const (
	tournamentKeyPrefix = "tournament:"
	tournamentListKey   = "tournaments"
	defaultExpiration   = 24 * time.Hour
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Create(ctx context.Context, tournament *pb.Tournament) error {
	tournamentJSON, err := json.Marshal(tournament)
	if err != nil {
		return fmt.Errorf("failed to marshal tournament: %w", err)
	}

	key := r.getTournamentKey(tournament.Id)
	pipe := r.client.Pipeline()

	pipe.Set(ctx, key, tournamentJSON, defaultExpiration)

	pipe.SAdd(ctx, tournamentListKey, tournament.Id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to store tournament: %w", err)
	}

	return nil
}

func (r *RedisRepository) Get(ctx context.Context, id string) (*pb.Tournament, error) {
	key := r.getTournamentKey(id)

	tournamentJSON, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("tournament not found")
		}
		return nil, fmt.Errorf("failed to get tournament: %w", err)
	}

	var tournament pb.Tournament
	if err := json.Unmarshal(tournamentJSON, &tournament); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tournament: %w", err)
	}

	return &tournament, nil
}

func (r *RedisRepository) List(ctx context.Context, pageSize int32, pageToken string, status pb.TournamentStatus) ([]*pb.Tournament, string, error) {
	tournamentIDs, err := r.client.SMembers(ctx, tournamentListKey).Result()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get tournament list: %w", err)
	}

	startIdx := 0
	if pageToken != "" {
		for i, id := range tournamentIDs {
			if id == pageToken {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + int(pageSize)
	if endIdx > len(tournamentIDs) {
		endIdx = len(tournamentIDs)
	}

	var tournaments []*pb.Tournament
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, 0, endIdx-startIdx)

	for i := startIdx; i < endIdx; i++ {
		cmd := pipe.Get(ctx, r.getTournamentKey(tournamentIDs[i]))
		cmds = append(cmds, cmd)
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, "", fmt.Errorf("failed to execute pipeline: %w", err)
	}

	for _, cmd := range cmds {
		tournamentJSON, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, "", fmt.Errorf("failed to get tournament data: %w", err)
		}

		var tournament pb.Tournament
		if err := json.Unmarshal(tournamentJSON, &tournament); err != nil {
			continue
		}

		if status != pb.TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED && tournament.Status != status {
			continue
		}

		tournaments = append(tournaments, &tournament)
	}

	var nextPageToken string
	if endIdx < len(tournamentIDs) {
		nextPageToken = tournamentIDs[endIdx]
	}

	return tournaments, nextPageToken, nil
}

func (r *RedisRepository) Update(ctx context.Context, tournament *pb.Tournament) error {
	tournamentJSON, err := json.Marshal(tournament)

	key := r.getTournamentKey(tournament.Id)
	pipe := r.client.Pipeline()

	if err != nil {
		return fmt.Errorf("failed to marshal tournament: %w", err)
	}

	pipe.Set(ctx, key, tournamentJSON, defaultExpiration).Err()

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update tournament: %w", err)
	}

	return nil
}

func (r *RedisRepository) getTournamentKey(id string) string {
	return fmt.Sprintf("%s%s", tournamentKeyPrefix, id)
}

type MemoryRepository struct {
	tournaments map[string]*pb.Tournament
	lock        *sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tournaments: make(map[string]*pb.Tournament),
		lock:        &sync.Mutex{},
	}
}

func (r *MemoryRepository) Create(ctx context.Context, tournament *pb.Tournament) error {
	r.lock.Lock()
	r.tournaments[tournament.Id] = tournament
	r.lock.Unlock()
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (*pb.Tournament, error) {
	tournament := r.tournaments[id]
	if tournament == nil {
		return nil, fmt.Errorf("tournament not found")
	}
	return tournament, nil
}

func (r *MemoryRepository) Update(ctx context.Context, tournament *pb.Tournament) error {
	if tournament == nil {
		return fmt.Errorf("tournament is nil")
	}
	r.lock.Lock()
	r.tournaments[tournament.Id] = tournament
	r.lock.Unlock()
	return nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]*pb.Tournament, error) {
	list := make([]*pb.Tournament, 0, len(r.tournaments))
	for _, tournament := range r.tournaments {
		list = append(list, tournament)
	}
	return list, nil
}
