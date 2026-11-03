package mojang

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewClient(c *cache.Cache) *Client {
	return &Client{cache: c}
}

func (c *Client) GetProfile(ctx context.Context, username string) (*Profile, error) {
	dest := "https://api.minecraftservices.com/minecraft/profile/lookup/name/" + url.QueryEscape(username)

	var profile Profile
	err := c.cache.GetJSON(ctx, "mojang:"+dest, 15*time.Minute, &profile, func(ctx context.Context) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, dest, nil)
	})
	if status.Code(err) == codes.NotFound {
		return nil, status.Errorf(codes.NotFound, "player %s not found", username)
	}

	if err != nil {
		return nil, err
	}

	profile.ID = strings.ToLower(strings.ReplaceAll(profile.ID, "-", ""))
	return &profile, nil
}
