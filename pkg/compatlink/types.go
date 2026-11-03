package compatlink

import (
	"encoding/json"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
)

type Client struct {
	baseURL string
	cache   *cache.Cache
}

type networthRequest struct {
	Member json.RawMessage `json:"member"`
	Museum json.RawMessage `json:"museum,omitempty"`
	Bank   *float64        `json:"bank,omitempty"`
}

type farmingWeightRequest struct {
	Profile json.RawMessage `json:"profile"`
	UUID    string          `json:"uuid"`
}

type NetworthResponse struct {
	UnsoulboundNetworth float64 `json:"unsoulboundNetworth"`
	Networth            float64 `json:"networth"`
}

type FarmingWeightResponse struct {
	TotalWeight float64 `json:"totalWeight"`
}
