package hypixel

import (
	"context"

	pb "github.com/guild-link/hypixel-grpc/proto"
)

func (s *Server) GetDungeons(ctx context.Context, req *pb.SkyBlockRequest) (*pb.DungeonsResponse, error) {
	d, err := s.hypixel.GetDungeons(ctx, req.GetUsername(), req.GetProfile())
	if err != nil {
		return nil, err
	}

	classes := make([]*pb.DungeonClassLevel, len(d.Classes))
	for i, class := range d.Classes {
		classes[i] = &pb.DungeonClassLevel{
			DungeonClass: pb.DungeonClass(class.Class),
			Level:        class.Level,
		}
	}

	floors := make([]*pb.DungeonFloorStats, len(d.Floors))
	for i, floor := range d.Floors {
		floors[i] = &pb.DungeonFloorStats{
			Mode:         pb.DungeonMode(floor.Mode),
			Floor:        floor.Floor,
			Completions:  floor.Completions,
			PersonalBest: floor.PersonalBest,
		}
	}

	response := &pb.DungeonsResponse{
		Profile:            profile(d.Profile),
		SelectedClassLevel: d.SelectedClassLevel,
		CatacombsLevel:     d.CatacombsLevel,
		ClassAverage:       d.ClassAverage,
		SecretsFound:       d.SecretsFound,
		Classes:            classes,
		Floors:             floors,
	}
	if d.SelectedClass != nil {
		selectedClass := pb.DungeonClass(*d.SelectedClass)
		response.SelectedClass = &selectedClass
	}

	return response, nil
}
