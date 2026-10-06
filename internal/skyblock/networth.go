package skyblock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetNetworth(ctx context.Context, req *pb.SkyBlockRequest) (*pb.NetworthResponse, error) {
	player, data, rawProfile, err := s.hypixel.GetRawProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	profile, err := hypixel.ParseProfile(player, data)
	if err != nil {
		return nil, err
	}

	museum, err := s.hypixel.GetMuseum(ctx, profile.ID, profile.Mojang.ID)
	if err != nil {
		return nil, err
	}

	var members struct {
		Members map[string]json.RawMessage `json:"members"`
	}

	if err := json.Unmarshal(rawProfile, &members); err != nil {
		return nil, fmt.Errorf("decode Hypixel profile: %w", err)
	}

	uuid := profile.Mojang.ID
	result, err := s.compat.Networth(ctx, members.Members[uuid], museum, data.Banking.Balance, uuid, profile.ID)
	if err != nil {
		return nil, err
	}

	return &pb.NetworthResponse{
		Unsoulbound: result.UnsoulboundNetworth,
		Profile:     profile.ProfileProto(),
		User:        profile.UserProto(),
		Total:       result.Networth,
	}, nil
}
