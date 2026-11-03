package skyblock

import (
	"github.com/guild-link/hypixel-grpc/pkg/calculator"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

func slayers(member hypixel.MemberData) *pb.Slayers {
	response := &pb.Slayers{}
	for value := 1; value < len(pb.SlayerType_name); value++ {
		slayer := pb.SlayerType(value)
		name := enumName(slayer, "SLAYER_TYPE_")

		xp := member.Slayer.SlayerBosses[name].XP
		response.Slayers = append(response.Slayers, &pb.SlayerLevel{
			Level:      calculator.SlayerLevel(name, xp),
			SlayerType: slayer,
			Xp:         xp,
		})
		response.TotalXp += xp
	}

	return response
}
