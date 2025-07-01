package server

import (
	"context"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// shortenServicer defines the interface for single URL shortening operations.
// Implementations should handle the creation of short URLs from original URLs.
type shortenServicer interface {
	// ShortenURL creates a shortened version of the original URL.
	ShortenURL(context.Context, models.UserID, models.OrigURL) (*models.URLPair, error)
}

// errShortenConflict defines the interface for URL conflict errors.
// Implementations should indicate when a URL already exists in the system.
type errShortenConflict interface {
	error
	// IsErrConflict returns true if the error represents a URL conflict
	IsErrConflict() bool
}

// ShortenURL handles single URL shortening requests.
//
// It validates the input URL, delegates the shortening operation to the service,
// and returns the shortened URL. Handles conflict cases gracefully by returning
// the existing short URL when available.
func (s *ShortenerServer) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	uid, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	_, err := url.ParseRequestURI(req.OriginalUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad request")
	}

	pair, err := s.shorten.ShortenURL(ctx, uid, models.OrigURL(req.OriginalUrl))
	if err != nil {
		if e, ok := err.(errShortenConflict); ok && e.IsErrConflict() {
			return &pb.ShortenURLResponse{
				ShortUrl: s.baseAddr + "/" + string(pair.Short),
			}, nil
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.ShortenURLResponse{
		ShortUrl: s.baseAddr + "/" + string(pair.Short),
	}, nil
}
