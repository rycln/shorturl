package interceptors

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// authServicer defines the interface for authentication operations.
// Implementations should handle both JWT generation and parsing.
type authServicer interface {
	// NewJWTString generates a new JWT token for the given user ID.
	NewJWTString(models.UserID) (string, error)

	// ParseIDFromAuthHeader extracts user ID from a JWT authorization header.
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

// AuthInterceptor implements gRPC unary server interceptor for authentication.
// It handles both existing JWT validation and new user registration.
type AuthInterceptor struct {
	authService authServicer // Service handling JWT operations
}

// NewAuthInterceptor creates a new AuthInterceptor instance.
func NewAuthInterceptor(authService authServicer) *AuthInterceptor {
	return &AuthInterceptor{
		authService: authService,
	}
}

// Auth performs authentication/authorization for gRPC requests.
//
// The interceptor:
// 1. Checks for existing Bearer token in Authorization header
// 2. Validates token if present and extracts user ID
// 3. Generates new token for new users
// 4. Sets user ID in request context
// 5. Adds new token to response headers when created
func (i *AuthInterceptor) Auth(ctx context.Context) (context.Context, error) {
	var userID models.UserID

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
			header := authHeaders[0]
			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")
				uid, err := i.authService.ParseIDFromAuthHeader(token)
				if err != nil {
					logger.Log.Debug("auth interceptor", zap.Error(err))
				} else {
					userID = uid
				}
			}
		}
	}

	if userID == "" {
		userID = models.UserID(uuid.NewString())

		jwtString, err := i.authService.NewJWTString(userID)
		if err != nil {
			logger.Log.Debug("auth interceptor", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to generate token")
		}

		header := metadata.Pairs("authorization", "Bearer "+jwtString)
		err = grpc.SetHeader(ctx, header)
		if err != nil {
			return nil, status.Error(codes.Internal, "can't write the header")
		}
	}

	return context.WithValue(ctx, contextkeys.UserID, userID), nil
}
