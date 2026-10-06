package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetProfile(ctx context.Context, username, profileName string) (*SkyBlockProfile, error) {
	player, data, _, err := c.GetRawProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	return ParseProfile(player, data)
}

func (c *Client) GetRawProfile(ctx context.Context, username, profileName string) (*mojang.Profile, *ProfileData, json.RawMessage, error) {
	player, err := c.mojang.GetProfile(ctx, username)
	if err != nil {
		return nil, nil, nil, err
	}

	body, err := c.get(ctx, "/skyblock/profiles?uuid="+url.QueryEscape(player.ID))
	if err != nil {
		return nil, nil, nil, err
	}

	var data struct {
		Profiles []json.RawMessage `json:"profiles"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, nil, fmt.Errorf("decode Hypixel profiles response: %w", err)
	}

	if len(data.Profiles) == 0 {
		return nil, nil, nil, status.Errorf(codes.NotFound, "%s has no SkyBlock profiles", player.Name)
	}

	for _, rawProfile := range data.Profiles {
		var profile ProfileData
		if err := json.Unmarshal(rawProfile, &profile); err != nil {
			return nil, nil, nil, fmt.Errorf("decode Hypixel profile: %w", err)
		}

		if (profileName == "" && profile.Selected) || strings.EqualFold(profile.CuteName, profileName) {
			return player, &profile, rawProfile, nil
		}
	}

	if profileName == "" {
		return nil, nil, nil, status.Errorf(codes.NotFound, "%s has no selected profile", player.Name)
	}

	return nil, nil, nil, status.Errorf(codes.NotFound, "%s has no profile named %s", player.Name, profileName)
}

func ParseProfile(player *mojang.Profile, data *ProfileData) (*SkyBlockProfile, error) {
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
