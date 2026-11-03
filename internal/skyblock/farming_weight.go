package skyblock

import (
	"context"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
)

func (s *Server) farmingWeight(ctx context.Context, profile *hypixel.Profile) (*float64, error) {
	result, err := s.compat.FarmingWeight(ctx, profile.Raw, profile.Player.ID, profile.ID)
	if err != nil {
		return nil, err
	}

	return &result.TotalWeight, nil
}
