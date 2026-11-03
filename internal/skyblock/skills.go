package skyblock

import (
	"math"

	"github.com/guild-link/hypixel-grpc/pkg/calculator"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func skills(member hypixel.MemberData, collectionTiers map[string]int) *pb.Skills {
	experience := member.PlayerData.Experience
	perks := member.PlayerData.Perks

	extraLevels := map[pb.Skill]int{
		pb.Skill_SKILL_FORAGING: max(perks["agatha_power"]-3, 0) + max(perks["miria_power"]-3, 0),
		pb.Skill_SKILL_TAMING:   len(member.PetsData.PetCare.PetTypesSacrificed),
		pb.Skill_SKILL_FARMING:  member.JacobContest.Perks.FarmingLevelCap,
	}

	for _, collection := range []string{"FIG_LOG", "MANGROVE_LOG", "HELIX_LOG"} {
		if collectionTiers[collection] >= 9 {
			extraLevels[pb.Skill_SKILL_FORAGING]++
		}
	}

	var levels []float64
	response := &pb.Skills{SkyblockLevel: calculator.SkyBlockLevel(member.Leveling.Experience)}
	for value := 1; value < len(pb.Skill_name); value++ {
		skill := pb.Skill(value)
		xp := experience[skill.String()]

		level := calculator.SkillLevel(skill.String(), xp, extraLevels[skill])
		response.Skills = append(response.Skills, &pb.SkillLevel{Skill: skill, Level: level, Xp: xp})
		levels = append(levels, math.Floor(level))
	}

	response.SkillAverage = calculator.Average(levels...)
	return response
}
