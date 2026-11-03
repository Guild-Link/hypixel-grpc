package skyblock

import (
	"strconv"

	"github.com/guild-link/hypixel-grpc/pkg/calculator"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func dungeons(member hypixel.MemberData) *pb.Dungeons {
	d := member.Dungeons
	response := &pb.Dungeons{
		CatacombsLevel: calculator.CatacombsLevel(d.DungeonTypes["catacombs"].Experience),
		SecretsFound:   uint32(d.Secrets),
	}

	var classLevels []float64
	for value := 1; value < len(pb.DungeonClass_name); value++ {
		class := pb.DungeonClass(value)
		name := enumName(class, "DUNGEON_CLASS_")

		xp := d.Classes[name].Experience
		level := calculator.CatacombsLevel(xp)
		response.Classes = append(response.Classes, &pb.DungeonClassLevel{DungeonClass: class, Level: level})
		classLevels = append(classLevels, min(level, 50))

		if name == d.SelectedDungeonClass {
			response.SelectedClass = class
		}
	}

	response.ClassAverage = calculator.Average(classLevels...)

	for value := 1; value < len(pb.DungeonMode_name); value++ {
		mode := pb.DungeonMode(value)
		data := d.DungeonTypes[enumName(mode, "DUNGEON_MODE_")]

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

	return response
}
