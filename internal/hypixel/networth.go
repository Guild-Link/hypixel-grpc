package hypixel

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetNetworth(ctx context.Context, req *pb.SkyBlockRequest) (*pb.NetworthResponse, error) {
	player, data, rawProfile, err := s.GetRawProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	profile, err := mkProfile(player, data)
	if err != nil {
		return nil, err
	}

	museum, err := s.GetMuseum(ctx, profile.ID, profile.Mojang.ID)
	if err != nil {
		return nil, err
	}

	var profileData struct {
		Members map[string]json.RawMessage `json:"members"`
		Banking *struct {
			Balance *float64 `json:"balance"`
		} `json:"banking"`
	}

	if err := json.Unmarshal(rawProfile, &profileData); err != nil {
		return nil, fmt.Errorf("decode Hypixel profile: %w", err)
	}

	var bank *float64
	if profileData.Banking != nil {
		bank = profileData.Banking.Balance
	}

	uuid := profile.Mojang.ID
	result, err := s.compat.Networth(ctx, profileData.Members[uuid], museum, bank, uuid, profile.ID)
	if err != nil {
		return nil, err
	}

	return &pb.NetworthResponse{
		Unsoulbound: result.UnsoulboundNetworth,
		Total:       result.Networth,
		Profile:     profile.Proto(),
	}, nil
}
