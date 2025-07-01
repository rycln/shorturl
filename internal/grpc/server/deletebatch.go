package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// deletionProcessor defines the interface for asynchronous URL deletion operations.
// Implementations should handle queuing URLs for deletion rather than performing
// immediate deletion, as this is typically a background operation.
type deletionProcessor interface {
	// AddURLsIntoDeletionQueue queues URLs for asynchronous deletion.
	AddURLsIntoDeletionQueue(models.UserID, []models.ShortURL)
}

// DeleteUserURLs handles batch URL deletion requests.
//
// This endpoint accepts a list of short URLs to delete and queues them for
// asynchronous processing. The operation completes immediately after queuing,
// while actual deletion happens in the background. Requires authentication.
func (s *ShortenerServer) DeleteUserURLs(
	ctx context.Context,
	req *pb.DeleteUserURLsRequest,
) (*emptypb.Empty, error) {
	uid, err := s.auth.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication failed")
	}

	surls := make([]models.ShortURL, len(req.ShortUrls))
	for i, url := range req.ShortUrls {
		surls[i] = models.ShortURL(url)
	}

	s.delProc.AddURLsIntoDeletionQueue(uid, surls)

	return &emptypb.Empty{}, nil
}
