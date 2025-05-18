package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type authServicer interface {
	NewJWTString(models.UserID) (string, error)
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

// AuthMiddleware provides JWT-based authentication middleware.
//
// The middleware:
// 1. Extracts and validates JWT token from Authorization header
// 2. If valid:
//   - Extracts user ID from token claims
//   - Stores user ID in request context
//
// 3. If invalid/missing:
//   - Generates new user ID
//   - Creates new JWT token
//   - Sets both in response cookies
//
// Expected header format:
//
//	Authorization: Bearer <token>
type AuthMiddleware struct {
	authService authServicer
}

// NewAuthMiddleware creates new auth middleware instance.
func NewAuthMiddleware(authService authServicer) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// JWT returns the middleware function for chi/router.
func (m *AuthMiddleware) JWT(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		var userID models.UserID
		if header := r.Header.Get("Authorization"); header != "" {
			uid, err := m.authService.ParseIDFromAuthHeader(header)
			if err != nil {
				logger.Log.Debug("auth middleware", zap.Error(err))
			} else {
				userID = uid
			}
		}

		if userID == "" {
			userID = models.UserID(uuid.NewString())

			jwtString, err := m.authService.NewJWTString(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Log.Debug("auth middleware", zap.Error(err))
				return
			}

			cookie := &http.Cookie{
				Name:  "jwt",
				Value: jwtString,
			}
			http.SetCookie(w, cookie)
			w.Header().Set("Authorization", "Bearer "+jwtString)
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(auth)
}
