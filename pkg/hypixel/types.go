package hypixel

type Profile struct {
	ProfileID string            `json:"profile_id"`
	CuteName  string            `json:"cute_name"`
	Selected  bool              `json:"selected"`
	Members   map[string]Member `json:"members"`
}

type Member struct {
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
