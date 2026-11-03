package hypixel

import (
	"context"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetGarden(ctx context.Context, profileID string) (*Garden, error) {
	var data struct {
		Garden Garden `json:"garden"`
	}

	err := c.getJSON(ctx, "/skyblock/garden?profile="+url.QueryEscape(profileID), &data)
	if status.Code(err) == codes.NotFound {
		return &Garden{}, nil
	}

	if err != nil {
		return nil, err
	}

	return &data.Garden, nil
}
