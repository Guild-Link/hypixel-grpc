package compatlink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
)

func NewClient(c *cache.Cache, baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		cache:   c,
	}
}

func (c *Client) Networth(ctx context.Context, member, museum json.RawMessage, bank *float64, uuid, profileID string) (*NetworthResponse, error) {
	return post[NetworthResponse](ctx, c, "/networth", uuid+":"+profileID, networthRequest{
		Member: member,
		Museum: museum,
		Bank:   bank,
	})
}

func (c *Client) FarmingWeight(ctx context.Context, profile json.RawMessage, uuid, profileID string) (*FarmingWeightResponse, error) {
	return post[FarmingWeightResponse](ctx, c, "/farming-weight", uuid+":"+profileID, farmingWeightRequest{
		Profile: profile,
		UUID:    uuid,
	})
}

func post[T any](ctx context.Context, c *Client, path, key string, data any) (*T, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", "compatlink", path, key)

	var response T
	err := c.cache.GetJSON(ctx, cacheKey, 5*time.Minute, &response, func(ctx context.Context) (*http.Request, error) {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
	if err != nil {
		return nil, err
	}

	return &response, nil
}
