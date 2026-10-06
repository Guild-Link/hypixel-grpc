package hypixel

import (
	"context"
	"strconv"
	"strings"

	"github.com/guild-link/hypixel-grpc/pkg/common"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetDungeons(ctx context.Context, req *pb.SkyBlockRequest) (*pb.DungeonsResponse, error) {
	profile, err := s.GetProfile(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	dung := profile.Data.Dungeons
	response := &pb.DungeonsResponse{
		Profile:        profile.Proto(),
		CatacombsLevel: hypixel.CatacombsLevel(dung.DungeonTypes["catacombs"].Experience),
		SecretsFound:   dung.Secrets,
		Classes:        make([]*pb.DungeonClassLevel, 0, len(pb.DungeonClass_name)),
		Floors:         make([]*pb.DungeonFloorStats, 0, 15),
	}

	classTotal := 0.0
	for value := range len(pb.DungeonClass_name) {
		class := pb.DungeonClass(value)
		name := strings.ToLower(strings.TrimPrefix(class.String(), "DUNGEON_CLASS_"))
		level := hypixel.CatacombsLevel(dung.Classes[name].Experience)
		response.Classes = append(response.Classes, &pb.DungeonClassLevel{DungeonClass: class, Level: level})

		classTotal += level
		if name == dung.SelectedDungeonClass {
			response.SelectedClass = &class
			response.SelectedClassLevel = level
		}
	}

	response.ClassAverage = common.RoundToTwo(classTotal / 5)

	floorStats := func(mode pb.DungeonMode, dungeonType string, floor uint32) *pb.DungeonFloorStats {
		data := dung.DungeonTypes[dungeonType]
		key := strconv.FormatUint(uint64(floor), 10)
		stats := &pb.DungeonFloorStats{
			Mode:        mode,
			Floor:       floor,
			Completions: data.TierCompletions[key],
		}

		if personalBest, ok := data.FastestTimeSPlus[key]; ok {
			stats.PersonalBest = &personalBest
		}

		return stats
	}

	for floor := uint32(0); floor <= 7; floor++ {
		response.Floors = append(response.Floors, floorStats(pb.DungeonMode_DUNGEON_MODE_CATACOMBS, "catacombs", floor))
	}

	for floor := uint32(1); floor <= 7; floor++ {
		response.Floors = append(response.Floors, floorStats(pb.DungeonMode_DUNGEON_MODE_MASTER_CATACOMBS, "master_catacombs", floor))
	}

	return response, nil
}
