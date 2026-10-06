package skyblock

import (
	"context"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetFarming(ctx context.Context, req *pb.SkyBlockRequest) (*pb.FarmingResponse, error) {
	player, data, rawProfile, err := s.hypixel.GetRawProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	profile, err := hypixel.ParseProfile(player, data)
	if err != nil {
		return nil, err
	}

	weight, err := s.compat.FarmingWeight(ctx, rawProfile, player.ID, profile.ID)
	if err != nil {
		return nil, err
	}

	farmingXP := profile.Data.PlayerData.Experience.SkillFarming
	gardenXP := profile.Data.Garden.Experience

	return &pb.FarmingResponse{
		FarmingLevel: hypixel.FarmingLevel(farmingXP),
		GardenLevel:  hypixel.GardenLevel(gardenXP),
		Weight:       weight.TotalWeight,
		FarmingXp:    farmingXP,
		GardenXp:     gardenXP,
		Profile:      profile.ProfileProto(),
		User:         profile.UserProto(),
	}, nil
}
