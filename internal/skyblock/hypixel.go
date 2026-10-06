package skyblock

import (
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/compatlink"
	"github.com/guild-link/hypixel-grpc/pkg/hypixel"
	pb "github.com/guild-link/hypixel-grpc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSkyBlockServer

	hypixel *hypixel.Client
	compat  *compatlink.Client
}

func Register(r grpc.ServiceRegistrar, c *cache.Cache, apiKey, compatURL string) {
	pb.RegisterSkyBlockServer(r, &Server{
		compat:  compatlink.NewClient(c, compatURL),
		hypixel: hypixel.NewClient(c, apiKey),
	})
}
