package hypixel

import (
	"context"
	"strconv"

	"github.com/guild-link/hypixel-grpc/pkg/math"
)

func (c *Client) GetDungeons(ctx context.Context, username, profileName string) (*DungeonStats, error) {
	profile, err := c.GetProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	calcCata := func(xp float64) float64 {
		return math.CalcXPTable(xp, catacombsXPTable[:], 200_000_000)
	}

	dung := profile.Data.Dungeons
	cl := dung.Classes

	classTotal := 0.0
	var selectedClassLevel float64
	var selectedClass *DungeonClass
	classes := make([]DungeonClassLevel, 0, len(dungeonClassNames))

	for value, name := range dungeonClassNames {
		class := DungeonClass(value)
		level := calcCata(cl[name].Experience)
		classes = append(classes, DungeonClassLevel{Class: class, Level: level})

		classTotal += level
		if name == dung.SelectedDungeonClass {
			selectedClass = &class
			selectedClassLevel = level
		}
	}

	floorStats := func(mode DungeonMode, dungeonType string, floor uint32) DungeonFloorStats {
		data := dung.DungeonTypes[dungeonType]
		key := strconv.FormatUint(uint64(floor), 10)
		personalBest, hasPersonalBest := data.FastestTimeSPlus[key]
		var personalBestValue *float64
		if hasPersonalBest {
			personalBestValue = &personalBest
		}

		return DungeonFloorStats{
			Mode:         mode,
			Floor:        floor,
			PersonalBest: personalBestValue,
			Completions:  data.TierCompletions[key],
		}
	}

	floors := make([]DungeonFloorStats, 0, 15)
	for floor := uint32(0); floor <= 7; floor++ {
		floors = append(floors, floorStats(DungeonModeCatacombs, "catacombs", floor))
	}

	for floor := uint32(1); floor <= 7; floor++ {
		floors = append(floors, floorStats(DungeonModeMasterCatacombs, "master_catacombs", floor))
	}

	return &DungeonStats{
		CatacombsLevel:     calcCata(dung.DungeonTypes["catacombs"].Experience),
		ClassAverage:       math.RoundToTwo(classTotal / 5),
		SelectedClassLevel: selectedClassLevel,
		SelectedClass:      selectedClass,
		SecretsFound:       dung.Secrets,
		Profile:            profile,
		Classes:            classes,
		Floors:             floors,
	}, nil
}
