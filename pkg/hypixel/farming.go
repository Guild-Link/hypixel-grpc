package hypixel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/guild-link/hypixel-grpc/pkg/math"
)

func (c *Client) GetFarming(ctx context.Context, username, profileName string) (*FarmingData, error) {
	player, rawProfile, err := c.GetRawProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	profile, err := parseRawProfile(player, rawProfile)
	if err != nil {
		return nil, err
	}
	weight, err := c.compat.FarmingWeight(ctx, rawProfile, player.ID, profile.ID)
	if err != nil {
		return nil, err
	}

	var rawData struct {
		Members map[string]struct {
			Garden struct {
				Experience float64 `json:"garden_experience"`
			} `json:"garden_player_data"`
		} `json:"members"`
	}

	if err := json.Unmarshal(rawProfile, &rawData); err != nil {
		return nil, fmt.Errorf("decode garden data: %w", err)
	}
	gardenXP := rawData.Members[player.ID].Garden.Experience

	var farmingXP float64
	if profile.Data.PlayerData != nil && profile.Data.PlayerData.Experience != nil {
		farmingXP = profile.Data.PlayerData.Experience.SkillFarming
	}

	return &FarmingData{
		GardenLevel:  math.CalcXPTable(gardenXP, gardenXPTable[:], 10000),
		FarmingLevel: math.CalcXPTable(farmingXP, farmingXPTable[:]),
		FarmingXP:    farmingXP,
		GardenXP:     gardenXP,
		Weight:       weight.TotalWeight,
	}, nil
}
