package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// retrieveServicer defines the interface for URL retrieval operations.
// Implementations should handle both URL lookups and context-based operations.
type retrieveServicer interface {
	// GetShortURLFromCtx retrieves a short URL from context.
	GetShortURLFromCtx(context.Context) (models.ShortURL, error)

	// GetOrigURLByShort looks up the original URL by its short version.
	GetOrigURLByShort(context.Context, models.ShortURL) (models.OrigURL, error)
}

// errRetrieveDeletedURL defines the interface for deleted URL errors.
// Implementations should indicate when a requested URL has been deleted.
type errRetrieveDeletedURL interface {
	error
	// IsErrDeletedURL returns true if the error represents a deleted URL
	IsErrDeletedURL() bool
}

// RetrieveURL handles URL retrieval requests by short URL identifier.
//
// It looks up the original URL corresponding to the provided short URL,
// handling special cases for deleted URLs and other error conditions.
func (s *ShortenerServer) RetrieveURL(
	ctx context.Context,
	req *pb.RetrieveURLRequest,
) (*pb.RetrieveURLResponse, error) {
	shortURL := models.ShortURL(req.ShortUrl)

	origURL, err := s.retrieve.GetOrigURLByShort(ctx, shortURL)
	if err != nil {
		if e, ok := err.(errRetrieveDeletedURL); ok && e.IsErrDeletedURL() {
			return nil, status.Error(codes.NotFound, "URL was deleted")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve URL")
	}

	return &pb.RetrieveURLResponse{
		OriginalUrl: string(origURL),
	}, nil
}
