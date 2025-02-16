package players

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "shared/grpc"

	"github.com/redis/go-redis/v9"
)

const (
	playerKeyPrefix   = "player:"
	playerListKey     = "players"
	defaultExpiration = 24 * time.Hour
)

type PlayersRepository struct {
	client *redis.Client
}

func NewPlayersRepository(client *redis.Client) *PlayersRepository {
	return &PlayersRepository{
		client: client,
	}
}

func (r *PlayersRepository) Create(ctx context.Context, player *pb.PlayerFile) error {
	playerJSON, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	key := r.getplayerKey(player.Id)
	pipe := r.client.Pipeline()

	pipe.Set(ctx, key, playerJSON, defaultExpiration)

	pipe.SAdd(ctx, playerKeyPrefix, player.Id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to store player: %w", err)
	}

	return nil
}

func (r *PlayersRepository) Get(ctx context.Context, id string) (*pb.PlayerFile, error) {
	key := r.getplayerKey(id)

	playerJSON, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	var player pb.PlayerFile
	if err := json.Unmarshal(playerJSON, &player); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player: %w", err)
	}

	return &player, nil
}

func (r *PlayersRepository) getplayerKey(id string) string {
	return fmt.Sprintf("%s%s", playerKeyPrefix, id)
}
