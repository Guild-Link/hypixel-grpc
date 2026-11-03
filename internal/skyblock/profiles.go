package skyblock

import (
	"context"
	"slices"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
	"golang.org/x/sync/errgroup"
)

var gameModes = map[string]pb.GameMode{
	"":        pb.GameMode_GAME_MODE_NORMAL,
	"bingo":   pb.GameMode_GAME_MODE_BINGO,
	"ironman": pb.GameMode_GAME_MODE_IRONMAN,
	"island":  pb.GameMode_GAME_MODE_STRANDED,
}

func (s *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	profile, err := s.hypixel.GetProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	result, err := s.skyBlockProfile(ctx, profile, req.GetInclude())
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileResponse{User: userProto(profile), Profile: result}, nil
}

func (s *Server) ListProfiles(ctx context.Context, req *pb.ListProfilesRequest) (*pb.ListProfilesResponse, error) {
	profiles, err := s.hypixel.GetProfiles(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	results := make([]*pb.SkyBlockProfile, len(profiles))
	eg, ctx := errgroup.WithContext(ctx)
	for i, profile := range profiles {
		eg.Go(func() (err error) {
			results[i], err = s.skyBlockProfile(ctx, profile, req.GetInclude())
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &pb.ListProfilesResponse{User: userProto(profiles[0]), Profiles: results}, nil
}

func (s *Server) skyBlockProfile(ctx context.Context, profile *hypixel.Profile, include []pb.Include) (*pb.SkyBlockProfile, error) {
	member := profile.Member
	result := &pb.SkyBlockProfile{
		Skills:   skills(member, profile.CollectionTiers),
		GameMode: gameModes[profile.GameMode],
		Purse:    member.Currencies.CoinPurse,
		Bank:     profile.Banking.Balance,
		Dungeons: dungeons(member),
		Selected: profile.Selected,
		Slayers:  slayers(member),
		Name:     profile.Name,
		Id:       profile.ID,
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, include := range slices.Compact(slices.Sorted(slices.Values(include))) {
		g.Go(func() (err error) {
			switch include {
			case pb.Include_INCLUDE_NETWORTH:
				result.Networth, err = s.networth(ctx, profile)
			case pb.Include_INCLUDE_GARDEN:
				result.Garden, err = s.garden(ctx, profile)
			case pb.Include_INCLUDE_FARMING_WEIGHT:
				result.FarmingWeight, err = s.farmingWeight(ctx, profile)
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
