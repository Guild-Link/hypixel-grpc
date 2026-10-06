package compatlink

import (
	"encoding/json"
	"net/http"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
)

type Client struct {
	baseURL string
	http    http.Client
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
	UnsoulboundNetworth float64      `json:"unsoulboundNetworth"`
	IsNonCosmetic       bool         `json:"isNonCosmetic"`
	PersonalBank        float64      `json:"personalBank"`
	NoInventory         bool         `json:"noInventory"`
	Networth            float64      `json:"networth"`
	Types               InventoryMap `json:"types"`
	Purse               float64      `json:"purse"`
	Bank                float64      `json:"bank"`
}

type InventoryMap map[string]NetworthInventory

type NetworthInventory struct {
	UnsoulboundTotal float64 `json:"unsoulboundTotal"`
	Total            float64 `json:"total"`
}

type FarmingWeightResponse struct {
	UncountedCrops map[string]float64 `json:"uncountedCrops"`
	BonusSources   map[string]float64 `json:"bonusSources"`
	TotalWeight    float64            `json:"totalWeight"`
	BonusWeight    float64            `json:"bonusWeight"`
	CropWeight     float64            `json:"cropWeight"`
}
