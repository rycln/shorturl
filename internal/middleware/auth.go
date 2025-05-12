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

type authServicer interface {
	NewJWTString(models.UserID) (string, error)
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type AuthMiddleware struct {
	authService authServicer
}

func NewAuthMiddleware(authService authServicer) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

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
		}

		jwtString, err := m.authService.NewJWTString(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Debug("auth middleware", zap.Error(err))
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userID)

		w.Header().Set("Authorization", "Bearer "+jwtString)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(auth)
}
