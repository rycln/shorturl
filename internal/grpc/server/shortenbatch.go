package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// shortenBatchServicer defines the interface for batch URL shortening operations.
// Implementations should handle the atomic processing of multiple URLs and return
// the complete set of shortened URLs.
type shortenBatchServicer interface {
	// BatchShortenURL processes multiple URLs in a single atomic operation.
	BatchShortenURL(context.Context, models.UserID, []models.OrigURL) ([]models.URLPair, error)
}

// BatchShortenURL handles batch URL shortening requests.
//
// It processes multiple URLs in a single operation while maintaining correlation
// between input and output items. The operation is atomic - either all URLs are
// shortened successfully or none are.
func (s *ShortenerServer) BatchShortenURL(
	ctx context.Context,
	req *pb.BatchShortenURLRequest,
) (*pb.BatchShortenURLResponse, error) {
	uid, ok := ctx.Value("userID").(models.UserID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	origs := make([]models.OrigURL, len(req.Items))
	for i, item := range req.Items {
		origs[i] = models.OrigURL(item.OriginalUrl)
	}

	pairs, err := s.batchShorten.BatchShortenURL(ctx, uid, origs)
	if err != nil {
		return nil, status.Error(codes.Internal, "batch processing failed")
	}

	res := &pb.BatchShortenURLResponse{
		Items: make([]*pb.BatchResultItem, len(pairs)),
	}

	for i, pair := range pairs {
		res.Items[i] = &pb.BatchResultItem{
			CorrelationId: req.Items[i].CorrelationId,
			ShortUrl:      s.baseAddr + "/" + string(pair.Short),
		}
	}

	return res, nil
}
