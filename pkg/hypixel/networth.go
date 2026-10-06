package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetNetworth(ctx context.Context, username, profileName string) (*Networth, error) {
	player, rawProfile, err := c.GetRawProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	profile, err := parseRawProfile(player, rawProfile)
	if err != nil {
		return nil, err
	}

	museum, err := c.getRawMuseum(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	var profileData struct {
		Members map[string]json.RawMessage `json:"members"`
		Banking *struct {
			Balance *float64 `json:"balance"`
		} `json:"banking"`
	}

	if err := json.Unmarshal(rawProfile, &profileData); err != nil {
		return nil, fmt.Errorf("decode Hypixel profile: %w", err)
	}

	var museumData struct {
		Members map[string]json.RawMessage `json:"members"`
	}

	if err := json.Unmarshal(museum, &museumData); err != nil {
		return nil, fmt.Errorf("decode Hypixel museum response: %w", err)
	}

	var bank *float64
	if profileData.Banking != nil {
		bank = profileData.Banking.Balance
	}

	uuid := profile.Mojang.ID
	result, err := c.compat.Networth(ctx, profileData.Members[uuid], museumData.Members[uuid], bank, uuid, profile.ID)
	if err != nil {
		return nil, err
	}

	return &Networth{
		Unsoulbound: result.UnsoulboundNetworth,
		Total:       result.Networth,
		Profile:     profile,
	}, nil
}
