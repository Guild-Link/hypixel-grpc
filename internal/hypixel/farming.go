package hypixel

import (
	"context"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetFarming(ctx context.Context, req *pb.SkyBlockRequest) (*pb.FarmingResponse, error) {
	player, data, rawProfile, err := s.GetRawProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	profile, err := mkProfile(player, data)
	if err != nil {
		return nil, err
	}

	weight, err := s.compat.FarmingWeight(ctx, rawProfile, player.ID, profile.ID)
	if err != nil {
		return nil, err
	}

	gardenXP := profile.Data.Garden.Experience
	farmingXP := profile.Data.PlayerData.Experience.SkillFarming

	return &pb.FarmingResponse{
		GardenLevel:  hypixel.GardenLevel(gardenXP),
		FarmingLevel: hypixel.FarmingLevel(farmingXP),
		FarmingXp:    farmingXP,
		GardenXp:     gardenXP,
		Weight:       weight.TotalWeight,
	}, nil
}
