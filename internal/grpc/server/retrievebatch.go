package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// retrieveBatchServicer defines the interface for batch URL retrieval operations.
// Implementations should handle fetching all URLs associated with a specific user.
type retrieveBatchServicer interface {
	// GetUserURLs retrieves all URL pairs (original and short) for a given user.
	GetUserURLs(context.Context, models.UserID) ([]models.URLPair, error)
}

// errRetrieveBatchNotExist defines the interface for non-existent user errors.
// Implementations should indicate when no URLs exist for the requested user.
type errRetrieveBatchNotExist interface {
	error
	// IsErrNotExist returns true if the error represents a "not found" condition
	IsErrNotExist() bool
}

// GetUserURLs retrieves all URLs belonging to the authenticated user.
//
// This endpoint requires authentication and returns all URL pairs (original and shortened)
// that the user has previously created. Returns an empty list if no URLs exist.
func (s *ShortenerServer) GetUserURLs(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.GetUserURLsResponse, error) {
	uid, err := s.auth.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication failed")
	}

	pairs, err := s.batchRetrieve.GetUserURLs(ctx, uid)
	if err != nil {
		if e, ok := err.(errRetrieveBatchNotExist); ok && e.IsErrNotExist() {
			return &pb.GetUserURLsResponse{Urls: nil}, nil
		}
		return nil, status.Error(codes.Internal, "failed to get user URLs")
	}

	res := &pb.GetUserURLsResponse{
		Urls: make([]*pb.UserURLItem, len(pairs)),
	}
	for i, pair := range pairs {
		res.Urls[i] = &pb.UserURLItem{
			ShortUrl:    s.baseAddr + "/" + string(pair.Short),
			OriginalUrl: string(pair.Orig),
		}
	}

	return res, nil
}
