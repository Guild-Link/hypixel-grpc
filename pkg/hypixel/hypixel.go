package hypixel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/compatlink"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
)

func NewClient(c *cache.Cache, apiKey, compatURL string) *Client {
	return &Client{
		cache:  c,
		apiKey: apiKey,
		http:   http.Client{Timeout: 15 * time.Second},
		mojang: mojang.NewClient(c),
		compat: compatlink.NewClient(compatURL),
	}
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	dest := fmt.Sprintf("%s/%s", "https://api.hypixel.net/v2", strings.TrimPrefix(path, "/"))
	cacheKey := fmt.Sprintf("%s:%s", "hypixel", dest)

	return c.cache.Do(ctx, cacheKey, func(ctx context.Context) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dest, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("API-Key", c.apiKey)

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("hypixel returned %s: %s", resp.Status, body)
		}

		return body, nil
	})
}
