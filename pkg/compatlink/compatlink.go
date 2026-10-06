package compatlink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewClient(c *cache.Cache, baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    http.Client{Timeout: 15 * time.Second},
		cache:   c,
	}
}

func (c *Client) Networth(ctx context.Context, member, museum json.RawMessage, bank *float64, uuid, profileID string) (*NetworthResponse, error) {
	response := &NetworthResponse{}
	err := c.post(ctx, "/networth", uuid+":"+profileID, networthRequest{
		Member: member,
		Museum: museum,
		Bank:   bank,
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) FarmingWeight(ctx context.Context, profile json.RawMessage, uuid, profileID string) (*FarmingWeightResponse, error) {
	response := &FarmingWeightResponse{}
	err := c.post(ctx, "/farming-weight", uuid+":"+profileID, farmingWeightRequest{
		Profile: profile,
		UUID:    uuid,
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) post(ctx context.Context, path, key string, data, response any) error {
	cacheKey := fmt.Sprintf("%s:%s:%s", "compatlink", path, key)
	respBody, err := c.cache.Do(ctx, cacheKey, 5*time.Minute, func(ctx context.Context) ([]byte, error) {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("encode compatlink request: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("create compatlink request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "compatlink unavailable: %v", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read compatlink response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, status.Errorf(codes.Unavailable, "compatlink POST %s returned %s: %s", path, resp.Status, respBody)
		}

		return respBody, nil
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, response); err != nil {
		return fmt.Errorf("decode compatlink response: %w", err)
	}

	return nil
}
