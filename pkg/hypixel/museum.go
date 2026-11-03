package hypixel

import (
	"context"
	"encoding/json"
	"net/url"
)

func (c *Client) GetMuseum(ctx context.Context, profileID, memberID string) (json.RawMessage, error) {
	var data struct {
		Members map[string]json.RawMessage `json:"members"`
	}

	if err := c.getJSON(ctx, "/skyblock/museum?profile="+url.QueryEscape(profileID), &data); err != nil {
		return nil, err
	}

	return data.Members[memberID], nil
}
