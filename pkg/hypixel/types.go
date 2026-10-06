package hypixel

import (
	"net/http"

	sc "github.com/DuckySoLucky/SkyCrypt-Types"
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/compatlink"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
)

type Client struct {
	apiKey string
	http   http.Client

	cache  *cache.Cache
	mojang *mojang.Client
	compat *compatlink.Client
}

type SkyBlockProfile struct {
	Mojang *mojang.Profile

	ID   string
	Name string

	Data *sc.Member
}

type Networth struct {
	Total       float64
	Unsoulbound float64
	Profile     *SkyBlockProfile
}
