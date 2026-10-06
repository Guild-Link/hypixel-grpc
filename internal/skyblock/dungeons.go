package skyblock

import (
	"context"
	"strconv"
	"strings"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetDungeons(ctx context.Context, req *pb.SkyBlockRequest) (*pb.DungeonsResponse, error) {
	profile, err := s.hypixel.GetProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	dung := profile.Data.Dungeons
	response := &pb.DungeonsResponse{
		CatacombsLevel: hypixel.CatacombsLevel(dung.DungeonTypes["catacombs"].Experience),
		SecretsFound:   uint32(dung.Secrets),
		Profile:        profile.ProfileProto(),
		User:           profile.UserProto(),
	}

	for value := 1; value < len(pb.DungeonClass_name); value++ {
		class := pb.DungeonClass(value)
		name := strings.ToLower(strings.TrimPrefix(class.String(), "DUNGEON_CLASS_"))

		level := hypixel.CatacombsLevel(dung.Classes[name].Experience)
		response.Classes = append(response.Classes, &pb.DungeonClassLevel{DungeonClass: class, Level: level})

		if name == dung.SelectedDungeonClass {
			response.SelectedClassLevel = level
			response.SelectedClass = class
		}
	}

	for value := 1; value < len(pb.DungeonMode_name); value++ {
		mode := pb.DungeonMode(value)
		data := dung.DungeonTypes[strings.ToLower(strings.TrimPrefix(mode.String(), "DUNGEON_MODE_"))]

		first := uint32(0)
		if mode == pb.DungeonMode_DUNGEON_MODE_MASTER_CATACOMBS {
			first = 1
		}

		for floor := first; floor <= 7; floor++ {
			key := strconv.FormatUint(uint64(floor), 10)
			stats := &pb.DungeonFloorStats{
				Completions: uint32(data.TierCompletions[key]),
				Floor:       floor,
				Mode:        mode,
			}

			if personalBest, ok := data.FastestTimeSPlus[key]; ok {
				stats.PersonalBest = &personalBest
			}

			response.Floors = append(response.Floors, stats)
		}
	}

	return response, nil
}
