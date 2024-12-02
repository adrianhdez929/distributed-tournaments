package games

import (
	"context"
	"fmt"
	"os"
	pb "shared/grpc"
)

type GamesService struct {
	pb.UnimplementedGamesServiceServer
	repo GamesRepository
}

func NewGamesService(repo GamesRepository) *GamesService {
	return &GamesService{
		repo: repo,
	}
}

func (s *GamesService) UploadGameFile(ctx context.Context, req *pb.UploadGameFileRequest) (*pb.UploadGameFileResponse, error) {
	// Validate the request
	if len(req.FileContent) == 0 || req.GameName == "" || req.ConstructorName == "" {
		return &pb.UploadGameFileResponse{
			Success: false,
		}, fmt.Errorf("invalid request")
	}

	// Save the file to disk (or handle it as needed)
	filePath := fmt.Sprintf("/uploads/%s.go", req.GameName)
	err := os.WriteFile(filePath, req.FileContent, 0644)

	if err != nil {
		return &pb.UploadGameFileResponse{
			Success: false,
		}, fmt.Errorf("failed to save file: %w", err)
	}

	return &pb.UploadGameFileResponse{
		Success: true,
	}, nil
}
