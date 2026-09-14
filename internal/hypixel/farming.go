package hypixel

import (
	"context"

	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetFarming(ctx context.Context, req *pb.SkyBlockRequest) (*pb.FarmingResponse, error) {
	farming, err := s.hypixel.GetFarming(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	return &pb.FarmingResponse{
		FarmingLevel: farming.FarmingLevel,
		GardenLevel:  farming.GardenLevel,
		FarmingXp:    farming.FarmingXP,
		GardenXp:     farming.GardenXP,
		Weight:       farming.Weight,
	}, nil
}
