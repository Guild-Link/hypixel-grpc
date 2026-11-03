package main

import (
	"cmp"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guild-link/hypixel-grpc/internal/skyblock"
	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cache, err := cache.NewCache(mustEnv("VALKEY_URL"), 30*time.Second)
	if err != nil {
		log.Fatalf("init cache: %v", err)
	}

	listener, err := net.Listen("tcp", cmp.Or(os.Getenv("LISTEN_ADDR"), ":50051"))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	skyblock.Register(server, cache, mustEnv("API_KEY"), mustEnv("COMPATLINK_URL"))

	context.AfterFunc(ctx, server.GracefulStop)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func mustEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s is not set", name)
	}

	return value
}
