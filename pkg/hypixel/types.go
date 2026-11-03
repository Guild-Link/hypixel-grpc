package hypixel

import (
	"encoding/json"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
)

type Client struct {
	apiKey string

	cache  *cache.Cache
	mojang *mojang.Client
}

type Profile struct {
	ID       string `json:"profile_id"`
	Name     string `json:"cute_name"`
	GameMode string `json:"game_mode"`
	Selected bool   `json:"selected"`

	CollectionTiers map[string]int `json:"-"`

	Raw     json.RawMessage            `json:"-"`
	Player  *mojang.Profile            `json:"-"`
	Member  MemberData                 `json:"-"`
	Members map[string]json.RawMessage `json:"members"`

	Banking struct {
		Balance *float64 `json:"balance"`
	} `json:"banking"`
}

type MemberData struct {
	Currencies struct {
		CoinPurse float64 `json:"coin_purse"`
	} `json:"currencies"`

	Leveling struct {
		Experience float64 `json:"experience"`
	} `json:"leveling"`

	PlayerData struct {
		UnlockedCollTiers []string           `json:"unlocked_coll_tiers"`
		Experience        map[string]float64 `json:"experience"`
		Perks             map[string]int     `json:"perks"`
	} `json:"player_data"`

	PetsData struct {
		PetCare struct {
			PetTypesSacrificed []string `json:"pet_types_sacrificed"`
		} `json:"pet_care"`
	} `json:"pets_data"`

	JacobContest struct {
		Perks struct {
			FarmingLevelCap int `json:"farming_level_cap"`
		} `json:"perks"`
	} `json:"jacobs_contest"`

	Slayer struct {
		SlayerBosses map[string]struct {
			XP float64 `json:"xp"`
		} `json:"slayer_bosses"`
	} `json:"slayer"`

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

type Garden struct {
	Experience     float64  `json:"garden_experience"`
	UnlockedPlots  []string `json:"unlocked_plots_ids"`
	CommissionData struct {
		TotalCompleted   float64 `json:"total_completed"`
		UniqueNPCsServed float64 `json:"unique_npcs_served"`
	} `json:"commission_data"`
}
