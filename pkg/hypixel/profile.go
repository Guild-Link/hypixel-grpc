package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	sc "github.com/DuckySoLucky/SkyCrypt-Types"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if len(rawProfiles) == 0 {
		return nil, nil, status.Errorf(codes.NotFound, "%s has no SkyBlock profiles", player.Name)
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
		return nil, nil, status.Errorf(codes.NotFound, "%s has no selected profile", player.Name)
	}

	return nil, nil, status.Errorf(codes.NotFound, "%s has no profile named %s", player.Name, profileName)
}

func parseRawProfile(player *mojang.Profile, rawProfile json.RawMessage) (*SkyBlockProfile, error) {
	var data sc.Profile
	if err := json.Unmarshal(rawProfile, &data); err != nil {
		return nil, fmt.Errorf("decode Hypixel profile: %w", err)
	}

	member, ok := data.Members[player.ID]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "%s is not a member of profile %s", player.Name, data.CuteName)
	}

	return &SkyBlockProfile{
		ID:     data.ProfileID,
		Name:   data.CuteName,
		Data:   &member,
		Mojang: player,
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

	return parseRawProfile(player, rawProfile)
}
