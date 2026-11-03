package skyblock

import (
	"context"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) networth(ctx context.Context, profile *hypixel.Profile) (*pb.Networth, error) {
	uuid := profile.Player.ID
	museum, err := s.hypixel.GetMuseum(ctx, profile.ID, uuid)
	if err != nil {
		return nil, err
	}

	result, err := s.compat.Networth(ctx, profile.Members[uuid], museum, profile.Banking.Balance, uuid, profile.ID)
	if err != nil {
		return nil, err
	}

	return &pb.Networth{Unsoulbound: result.UnsoulboundNetworth, Total: result.Networth}, nil
}
