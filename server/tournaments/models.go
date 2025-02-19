package tournaments

import (
	"context"

	pb "shared/grpc"
)

type TournamentRepository interface {
	Create(ctx context.Context, tournament *pb.Tournament) error
	Get(ctx context.Context, id string) (*pb.Tournament, error)
	List(ctx context.Context) ([]*pb.Tournament, error)
	Update(ctx context.Context, tournament *pb.Tournament) error
}
