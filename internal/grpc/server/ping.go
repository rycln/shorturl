package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// pingServicer defines the interface for storage health check operations.
// Implementations should verify connectivity to the underlying data storage.
type pingServicer interface {
	// PingStorage checks the connection to the data storage layer.
	PingStorage(context.Context) error
}

// Ping performs a health check of the server's connection to its storage backend.
//
// This endpoint provides a way to verify the service's ability to communicate with
// its persistent storage. It returns successfully only when the storage is accessible.
func (s *ShortenerServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.ping.PingStorage(ctx); err != nil {
		return nil, status.Error(codes.Internal, "storage connection failed")
	}
	return &emptypb.Empty{}, nil
}
