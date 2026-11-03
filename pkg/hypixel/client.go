package hypixel

import (
	"context"
	"net/http"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
)

func NewClient(c *cache.Cache, apiKey string) *Client {
	return &Client{
		cache:  c,
		apiKey: apiKey,
		mojang: mojang.NewClient(c),
	}
}

func (c *Client) getJSON(ctx context.Context, path string, v any) error {
	dest := "https://api.hypixel.net/v2" + path
	return c.cache.GetJSON(ctx, "hypixel:"+dest, 15*time.Minute, v, func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dest, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("API-Key", c.apiKey)
		return req, nil
	})
}
