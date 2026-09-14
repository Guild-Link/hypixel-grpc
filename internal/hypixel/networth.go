package hypixel

import (
	"context"

	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetNetworth(ctx context.Context, req *pb.SkyBlockRequest) (*pb.NetworthResponse, error) {
	nw, err := s.hypixel.GetNetworth(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	return &pb.NetworthResponse{
		Total:       nw.Total,
		Unsoulbound: nw.Unsoulbound,
		Profile:     profile(nw.Profile),
	}, nil
}
