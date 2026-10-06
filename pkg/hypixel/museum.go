package hypixel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) GetMuseum(ctx context.Context, profileID, memberID string) (json.RawMessage, error) {
	body, err := c.get(ctx, "/skyblock/museum?profile="+url.QueryEscape(profileID))
	if err != nil {
		return nil, err
	}

	var data struct {
		Members map[string]json.RawMessage `json:"members"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("decode Hypixel museum response: %w", err)
	}

	return data.Members[memberID], nil
}
