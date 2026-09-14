package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guild-link/hypixel-grpc/internal/hypixel"
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/common"
	"google.golang.org/grpc"
)

func newCache(url string) (*cache.Cache, error) {
	if url == "" {
		return nil, nil
	}

	return cache.NewCache(url, 15*time.Minute, 30*time.Second)
}

func serve(ctx context.Context, reg func(s *grpc.Server)) error {
	listener, err := net.Listen("tcp", common.DefaultEnv("LISTEN_ADDR", ":50051"))
	if err != nil {
		return err
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	reg(grpcServer)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listener)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	compatURL := common.MustEnv("COMPATLINK_URL")
	apiKey := common.MustEnv("API_KEY")

	c, err := newCache(os.Getenv("VALKEY_URL"))
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}

	reg := func(s *grpc.Server) {
		hypixel.Register(s, c, apiKey, compatURL)
	}

	if err := serve(ctx, reg); err != nil {
		log.Fatal(err)
	}
}
