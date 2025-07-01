package server

import (
	"context"
	"net"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// statsServicer defines the interface for statistics retrieval operations.
// Implementations should handle fetching system statistics from persistent storage.
type statsServicer interface {
	// GetStatsFromStorage retrieves current system statistics.
	GetStatsFromStorage(context.Context) (*models.Stats, error)
}

// GetStats retrieves system statistics (URLs and users counts).
//
// This endpoint is restricted to clients from trusted subnets only. It provides
// operational metrics about the service's usage and adoption.
func (s *ShortenerServer) GetStats(ctx context.Context, _ *emptypb.Empty) (*pb.GetStatsResponse, error) {
	if err := s.checkTrustedSubnet(ctx); err != nil {
		return nil, err
	}

	stats, err := s.stats.GetStatsFromStorage(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get stats from storage")
	}

	return &pb.GetStatsResponse{
		Urls:  uint64(stats.URLs),
		Users: uint64(stats.Users),
	}, nil
}

// checkTrustedSubnet verifies if the client's IP is within the trusted subnet.
func (s *ShortenerServer) checkTrustedSubnet(ctx context.Context) error {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Error(codes.PermissionDenied, "peer info not available")
	}

	clientIP := parseIPFromAddr(p.Addr.String())

	_, subnet, err := net.ParseCIDR(s.trustedSubnet)
	if err != nil {
		return status.Error(codes.Internal, "invalid trusted subnet configuration")
	}

	if clientIP == nil || !subnet.Contains(clientIP) {
		return status.Error(codes.PermissionDenied, "access denied: untrusted network")
	}

	return nil
}

// parseIPFromAddr extracts the IP address from a network address string.
func parseIPFromAddr(addr string) net.IP {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return net.ParseIP(addr)
	}
	return net.ParseIP(host)
}
