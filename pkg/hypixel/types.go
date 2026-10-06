package hypixel

import (
	"net/http"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

type Client struct {
	apiKey string
	http   http.Client

	cache  *cache.Cache
	mojang *mojang.Client
}

type SkyBlockProfile struct {
	Mojang *mojang.Profile
	Data   *MemberData

	Name string
	ID   string
}

func (p *SkyBlockProfile) ProfileProto() *pb.Profile {
	return &pb.Profile{Name: p.Name, Id: p.ID}
}

func (p *SkyBlockProfile) UserProto() *pb.User {
	return &pb.User{Name: p.Mojang.Name, Id: p.Mojang.ID}
}

type ProfileData struct {
	ProfileID string                `json:"profile_id"`
	CuteName  string                `json:"cute_name"`
	Selected  bool                  `json:"selected"`
	Members   map[string]MemberData `json:"members"`

	Banking struct {
		Balance *float64 `json:"balance"`
	} `json:"banking"`
}

type MemberData struct {
	PlayerData struct {
		Experience struct {
			SkillFarming float64 `json:"SKILL_FARMING"`
		} `json:"experience"`
	} `json:"player_data"`

	Garden struct {
		Experience float64 `json:"garden_experience"`
	} `json:"garden_player_data"`

	Dungeons struct {
		DungeonTypes map[string]struct {
			Experience       float64            `json:"experience"`
			TierCompletions  map[string]float64 `json:"tier_completions"`
			FastestTimeSPlus map[string]float64 `json:"fastest_time_s_plus"`
		} `json:"dungeon_types"`

		Classes map[string]struct {
			Experience float64 `json:"experience"`
		} `json:"player_classes"`

		SelectedDungeonClass string  `json:"selected_dungeon_class"`
		Secrets              float64 `json:"secrets"`
	} `json:"dungeons"`
}
