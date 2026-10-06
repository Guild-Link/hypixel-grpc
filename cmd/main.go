package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/guild-link/hypixel-grpc/internal/skyblock"
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	compatURL := common.MustEnv("COMPATLINK_URL")
	apiKey := common.MustEnv("API_KEY")

	cache, err := cache.NewCache(common.MustEnv("VALKEY_URL"), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}

	listener, err := net.Listen("tcp", common.DefaultEnv("LISTEN_ADDR", ":50051"))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if _, ok := status.FromError(err); !ok {
			err = status.Error(codes.Internal, err.Error())
		}

		return resp, err
	}))

	skyblock.Register(server, cache, apiKey, compatURL)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
