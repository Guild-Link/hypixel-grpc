package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyaside"
)

func NewCache(connectionString string, timeout time.Duration) (*Cache, error) {
	options, err := valkey.ParseURL(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse Valkey URL: %w", err)
	}

	client, err := valkeyaside.NewClient(valkeyaside.ClientOption{ClientOption: options})
	if err != nil {
		return nil, fmt.Errorf("create valkey client: %w", err)
	}

	return &Cache{client: client, timeout: timeout}, nil
}

func (c *Cache) Do(ctx context.Context, key string, ttl time.Duration, fn func(context.Context) ([]byte, error)) ([]byte, error) {
	cacheCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	value, err := c.client.Get(cacheCtx, ttl, key, func(ctx context.Context, _ string) (string, error) {
		data, err := fn(ctx)
		return string(data), err
	})
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}
