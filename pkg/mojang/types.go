package mojang

import "github.com/guild-link/hypixel-grpc/pkg/cache"

type Client struct {
	cache *cache.Cache
}

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
