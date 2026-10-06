package hypixel

import (
	"net/http"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/compatlink"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	pb "github.com/guild-link/hypixel-grpc/proto"
)

type Server struct {
	pb.UnimplementedHypixelServer

	apiKey string
	http   http.Client

	cache  *cache.Cache
	mojang *mojang.Client
	compat *compatlink.Client
}

type SkyBlockProfile struct {
	Mojang *mojang.Profile
	Data   *hypixel.Member

	Name string
	ID   string
}

func (p *SkyBlockProfile) Proto() *pb.SkyBlockProfile {
	return &pb.SkyBlockProfile{
		Profile: p.Name,
		Mojang: &pb.MojangProfile{
			Username: p.Mojang.Name,
			Uuid:     p.Mojang.ID,
		},
	}
}
