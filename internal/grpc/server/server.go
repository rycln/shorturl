package server

import (
	"context"

	pb "github.com/rycln/shorturl/api/gen/shortener"
	"github.com/rycln/shorturl/internal/models"
)

// authServicer defines the interface for authentication-related operations.
// Implementations should handle user identification from context.
type authServicer interface {
	// GetUserIDFromCtx extracts the user ID from the request context.
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// ShortenerServer implements the gRPC ShortenerService server.
//
// It provides URL shortening, retrieval, and management functionality
// through dependency-injected service components. The server handles:
// - Single and batch URL shortening
// - URL retrieval (single and batch)
// - URL deletion
// - System health checks
// - Usage statistics
//
// The server requires a base address for constructing full short URLs
// and optionally a trusted subnet for admin functionality.
type ShortenerServer struct {
	pb.UnimplementedShortenerServiceServer

	shorten       shortenServicer       // Handles single URL shortening
	batchShorten  shortenBatchServicer  // Handles batch URL shortening
	retrieve      retrieveServicer      // Handles single URL retrieval
	batchRetrieve retrieveBatchServicer // Handles batch URL retrieval
	delProc       deletionProcessor     // Handles URL deletion processing
	auth          authServicer          // Handles user authentication
	ping          pingServicer          // Handles health checks
	stats         statsServicer         // Handles statistics collection
	baseAddr      string                // Base address for short URLs
	trustedSubnet string                // Trusted subnet (CIDR notation)
}

// NewShortenerServer creates and initializes a new ShortenerServer instance.
//
// This constructor injects all required dependencies and configures the server
// with necessary operational parameters.
func NewShortenerServer(
	shorten shortenServicer,
	batchShorten shortenBatchServicer,
	retrieve retrieveServicer,
	batchRetrieve retrieveBatchServicer,
	delProc deletionProcessor,
	auth authServicer,
	ping pingServicer,
	stats statsServicer,
	baseAddr string,
	trustedSubnet string,
) *ShortenerServer {
	return &ShortenerServer{
		shorten:       shorten,
		batchShorten:  batchShorten,
		retrieve:      retrieve,
		batchRetrieve: batchRetrieve,
		delProc:       delProc,
		auth:          auth,
		ping:          ping,
		stats:         stats,
		baseAddr:      baseAddr,
		trustedSubnet: trustedSubnet,
	}
}
