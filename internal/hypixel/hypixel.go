package hypixel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/guild-link/hypixel-grpc/pkg/cache"
	"github.com/guild-link/hypixel-grpc/pkg/compatlink"
	"github.com/guild-link/hypixel-grpc/pkg/mojang"
	pb "github.com/guild-link/hypixel-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Register(r grpc.ServiceRegistrar, c *cache.Cache, apiKey, compatURL string) {
	pb.RegisterHypixelServer(r, &Server{
		cache:  c,
		apiKey: apiKey,
		mojang: mojang.NewClient(c),
		compat: compatlink.NewClient(c, compatURL),
		http:   http.Client{Timeout: 15 * time.Second},
	})
}

func (s *Server) get(ctx context.Context, path string) ([]byte, error) {
	dest := fmt.Sprintf("%s/%s", "https://api.hypixel.net/v2", strings.TrimPrefix(path, "/"))
	cacheKey := fmt.Sprintf("%s:%s", "hypixel", dest)

	return s.cache.Do(ctx, cacheKey, 15*time.Minute, func(ctx context.Context) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dest, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("API-Key", s.apiKey)

		resp, err := s.http.Do(req)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "Hypixel API unavailable: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, status.Error(codes.ResourceExhausted, "Hypixel API rate limit reached, try again later")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, status.Errorf(codes.Unavailable, "Hypixel API returned %s: %s", resp.Status, body)
		}

		return body, nil
	})
}
