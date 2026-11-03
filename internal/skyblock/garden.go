package skyblock

import (
	"context"

	"github.com/guild-link/hypixel-grpc/pkg/calculator"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) garden(ctx context.Context, profile *hypixel.Profile) (*pb.Garden, error) {
	garden, err := s.hypixel.GetGarden(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	return &pb.Garden{
		UniqueVisitors: uint32(garden.CommissionData.UniqueNPCsServed),
		VisitorsServed: uint32(garden.CommissionData.TotalCompleted),
		PlotsUnlocked:  uint32(len(garden.UnlockedPlots)),
		Level:          calculator.GardenLevel(garden.Experience),
		Xp:             garden.Experience,
	}, nil
}
