package games

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "shared/grpc"

	"github.com/redis/go-redis/v9"
)

const (
	gameKeyPrefix     = "game:"
	gameListKey       = "games"
	defaultExpiration = 24 * time.Hour
)

type GamesRepository struct {
	client *redis.Client
}

func NewGamesRepository(client *redis.Client) *GamesRepository {
	return &GamesRepository{
		client: client,
	}
}

func (r *GamesRepository) Create(ctx context.Context, game *pb.GameFile) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game: %w", err)
	}

	key := r.getGameKey(game.Id)
	pipe := r.client.Pipeline()

	pipe.Set(ctx, key, gameJSON, defaultExpiration)

	pipe.SAdd(ctx, gameKeyPrefix, game.Id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to store game: %w", err)
	}

	return nil
}

func (r *GamesRepository) Get(ctx context.Context, id string) (*pb.GameFile, error) {
	key := r.getGameKey(id)

	gameJSON, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("game not found")
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	var game pb.GameFile
	if err := json.Unmarshal(gameJSON, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game: %w", err)
	}

	return &game, nil
}

func (r *GamesRepository) getGameKey(id string) string {
	return fmt.Sprintf("%s%s", gameKeyPrefix, id)
}
