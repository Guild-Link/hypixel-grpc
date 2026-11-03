package hypixel

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetProfile(ctx context.Context, username, profileName string) (*Profile, error) {
	profiles, err := c.GetProfiles(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		if (profileName == "" && profile.Selected) || strings.EqualFold(profile.Name, profileName) {
			return profile, nil
		}
	}

	player := profiles[0].Player
	if profileName == "" {
		return nil, status.Errorf(codes.NotFound, "%s has no selected profile", player.Name)
	}

	return nil, status.Errorf(codes.NotFound, "%s has no profile named %s", player.Name, profileName)
}

func (c *Client) GetProfiles(ctx context.Context, username string) ([]*Profile, error) {
	player, err := c.mojang.GetProfile(ctx, username)
	if err != nil {
		return nil, err
	}

	var data struct {
		Profiles []json.RawMessage `json:"profiles"`
	}

	if err := c.getJSON(ctx, "/skyblock/profiles?uuid="+url.QueryEscape(player.ID), &data); err != nil {
		return nil, err
	}

	var profiles []*Profile
	for _, raw := range data.Profiles {
		profile := &Profile{Player: player, Raw: raw}
		if err := json.Unmarshal(raw, profile); err != nil {
			return nil, err
		}

		member, ok := profile.Members[player.ID]
		if !ok {
			continue
		}

		if err := json.Unmarshal(member, &profile.Member); err != nil {
			return nil, err
		}

		profile.CollectionTiers = map[string]int{}
		for _, member := range profile.Members {
			var data MemberData
			if err := json.Unmarshal(member, &data); err != nil {
				return nil, err
			}

			for _, unlocked := range data.PlayerData.UnlockedCollTiers {
				i := strings.LastIndex(unlocked, "_")
				tier, err := strconv.Atoi(unlocked[i+1:])
				if err != nil {
					continue
				}

				name := unlocked[:i]
				profile.CollectionTiers[name] = max(profile.CollectionTiers[name], tier)
			}
		}

		profiles = append(profiles, profile)
	}

	if len(profiles) == 0 {
		return nil, status.Errorf(codes.NotFound, "%s has no SkyBlock profiles", player.Name)
	}

	return profiles, nil
}
