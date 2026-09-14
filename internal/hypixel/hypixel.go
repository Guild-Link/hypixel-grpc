package hypixel

import (
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedHypixelServer
	hypixel *hypixel.Client
}

func profile(profile *hypixel.SkyBlockProfile) *pb.SkyBlockProfile {
	return &pb.SkyBlockProfile{
		Profile: profile.Name,
		Mojang: &pb.MojangProfile{
			Username: profile.Mojang.Name,
			Uuid:     profile.Mojang.ID,
		},
	}
}

func Register(s grpc.ServiceRegistrar, c *cache.Cache, apiKey, compatURL string) {
	pb.RegisterHypixelServer(s, &Server{
		hypixel: hypixel.NewClient(c, apiKey, compatURL),
	})
}
