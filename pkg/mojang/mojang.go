package mojang

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewClient(cache *cache.Cache) *Client {
	return &Client{
		cache: cache,
		http:  http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) GetProfile(ctx context.Context, username string) (*Profile, error) {
	dest := "https://api.minecraftservices.com/minecraft/profile/lookup/name/" + url.QueryEscape(username)
	cacheKey := fmt.Sprintf("%s:%s", "mojang", dest)

	body, err := c.cache.Do(ctx, cacheKey, 15*time.Minute, func(ctx context.Context) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dest, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "Mojang API unavailable: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, status.Errorf(codes.NotFound, "player %s not found", username)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, status.Errorf(codes.Unavailable, "Mojang API returned %s: %s", resp.Status, body)
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	profile.ID = strings.ToLower(strings.ReplaceAll(profile.ID, "-", ""))
	return &profile, nil
}
