package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	sc "github.com/DuckySoLucky/SkyCrypt-Types"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
)

func (c *Client) GetRawProfiles(ctx context.Context, username string) (*mojang.Profile, []json.RawMessage, error) {
	profile, err := c.mojang.GetProfile(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	body, err := c.get(ctx, "/skyblock/profiles?uuid="+url.QueryEscape(profile.ID))
	if err != nil {
		return nil, nil, err
	}

	var data struct {
		Profiles []json.RawMessage `json:"profiles"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, fmt.Errorf("decode Hypixel profiles response: %w", err)
	}

	return profile, data.Profiles, nil
}

func (c *Client) GetRawProfile(ctx context.Context, username, profileName string) (*mojang.Profile, json.RawMessage, error) {
	player, rawProfiles, err := c.GetRawProfiles(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	for _, rawProfile := range rawProfiles {
		var profile struct {
			Name     string `json:"cute_name"`
			Selected bool   `json:"selected"`
		}
		if err := json.Unmarshal(rawProfile, &profile); err != nil {
			return nil, nil, fmt.Errorf("decode Hypixel profile: %w", err)
		}

		if (profileName == "" && profile.Selected) || strings.EqualFold(profile.Name, profileName) {
			return player, rawProfile, nil
		}
	}

	if profileName == "" {
		return nil, nil, fmt.Errorf("selected profile not found")
	}

	return nil, nil, fmt.Errorf("profile %s not found", profileName)
}

func parseRawProfile(player *mojang.Profile, rawProfile json.RawMessage) (*SkyBlockProfile, error) {
	var data sc.Profile
	if err := json.Unmarshal(rawProfile, &data); err != nil {
		return nil, fmt.Errorf("decode Hypixel profile: %w", err)
	}

	var memberData *sc.Member
	if member, ok := data.Members[player.ID]; ok {
		memberData = &member
	}

	return &SkyBlockProfile{
		ID:                data.ProfileID,
		Name:              data.CuteName,
		Mojang:            player,
		GameMode:          data.GameMode,
		Selected:          data.Selected,
		Data:              memberData,
		Banking:           data.Banking,
		CommunityUpgrades: data.CommunityUpgrades,
	}, nil
}

func (c *Client) GetProfiles(ctx context.Context, username string) ([]SkyBlockProfile, error) {
	player, rawProfiles, err := c.GetRawProfiles(ctx, username)
	if err != nil {
		return nil, err
	}

	profiles := make([]SkyBlockProfile, 0, len(rawProfiles))
	for _, rawProfile := range rawProfiles {
		profile, err := parseRawProfile(player, rawProfile)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}

	return profiles, nil
}

func (c *Client) GetProfile(ctx context.Context, username, profileName string) (*SkyBlockProfile, error) {
	player, rawProfile, err := c.GetRawProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	profile, err := parseRawProfile(player, rawProfile)
	if err != nil {
		return nil, err
	}
	if profile.Data == nil {
		return nil, fmt.Errorf("member %s not found in profile %s", player.ID, profile.Name)
	}

	return profile, nil
}
