package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyaside"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var httpClient = http.Client{Timeout: 15 * time.Second}

func NewCache(valkeyURL string, timeout time.Duration) (*Cache, error) {
	options, err := valkey.ParseURL(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("parse Valkey URL: %w", err)
	}

	client, err := valkeyaside.NewClient(valkeyaside.ClientOption{ClientOption: options})
	if err != nil {
		return nil, fmt.Errorf("create Valkey client: %w", err)
	}

	return &Cache{client: client, timeout: timeout}, nil
}

func (c *Cache) GetJSON(ctx context.Context, key string, ttl time.Duration, v any, f func(context.Context) (*http.Request, error)) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	value, err := c.client.Get(ctx, ttl, key, func(ctx context.Context, _ string) (string, error) {
		req, err := f(ctx)
		if err != nil {
			return "", err
		}

		body, err := sendRequest(req)
		return string(body), err
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), v); err != nil {
		return fmt.Errorf("decode %s: %w", key, err)
	}

	return nil
}

func sendRequest(req *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "%s unavailable: %v", req.URL.Host, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", req.URL.Host, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.StatusNotFound:
		return nil, status.Errorf(codes.NotFound, "%s returned %s", req.URL.Host, resp.Status)
	case http.StatusTooManyRequests:
		return nil, status.Errorf(codes.ResourceExhausted, "%s rate limit reached, try again later", req.URL.Host)
	}

	return nil, status.Errorf(codes.Unavailable, "%s returned %s: %s", req.URL.Host, resp.Status, body)
}
