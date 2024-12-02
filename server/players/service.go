package players

import (
	"context"
	"fmt"
	"os"
	pb "shared/grpc"
)

type PlayersService struct {
	pb.UnimplementedPlayersServiceServer
	repo PlayersRepository
}

func NewPlayersService(repo PlayersRepository) *PlayersService {
	return &PlayersService{
		repo: repo,
	}
}

func (s *PlayersService) UploadGameFile(ctx context.Context, req *pb.UploadPlayerFileRequest) (*pb.UploadPlayerFileResponse, error) {
	// Validate the request
	if len(req.FileContent) == 0 || req.PlayerName == "" || req.ConstructorName == "" {
		return &pb.UploadPlayerFileResponse{
			Success: false,
		}, fmt.Errorf("invalid request")
	}

	// Save the file to disk (or handle it as needed)
	filePath := fmt.Sprintf("/uploads/%s.go", req.PlayerName)
	err := os.WriteFile(filePath, req.FileContent, 0644)

	if err != nil {
		return &pb.UploadPlayerFileResponse{
			Success: false,
		}, fmt.Errorf("failed to save file: %w", err)
	}

	return &pb.UploadPlayerFileResponse{
		Success: true,
	}, nil
}
