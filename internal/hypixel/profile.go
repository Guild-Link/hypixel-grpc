package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRawProfiles(ctx context.Context, username string) (*mojang.Profile, []json.RawMessage, error) {
	profile, err := s.mojang.GetProfile(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	body, err := s.get(ctx, "/skyblock/profiles?uuid="+url.QueryEscape(profile.ID))
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

func (s *Server) GetRawProfile(ctx context.Context, username, profileName string) (*mojang.Profile, *hypixel.Profile, json.RawMessage, error) {
	player, rawProfiles, err := s.GetRawProfiles(ctx, username)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(rawProfiles) == 0 {
		return nil, nil, nil, status.Errorf(codes.NotFound, "%s has no SkyBlock profiles", player.Name)
	}

	for _, rawProfile := range rawProfiles {
		var profile hypixel.Profile
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

func (s *Server) GetProfile(ctx context.Context, username, profileName string) (*SkyBlockProfile, error) {
	player, data, _, err := s.GetRawProfile(ctx, username, profileName)
	if err != nil {
		return nil, err
	}

	return mkProfile(player, data)
}

func mkProfile(player *mojang.Profile, data *hypixel.Profile) (*SkyBlockProfile, error) {
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
