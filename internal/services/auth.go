package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/models"
)

var errNoUserID = errors.New("does not contain user id")

// Auth provides user authentication services using JWT tokens.
//
// The service handles token generation, validation and user ID extraction
// from request context. It's designed to work seamlessly with HTTP middleware.
type Auth struct {
	jwtKey string
	jwtExp time.Duration
}

// NewAuth creates a new Auth service instance with configured secret and token expiration value.
func NewAuth(jwtkey string, jwtExp time.Duration) *Auth {
	return &Auth{
		jwtKey: jwtkey,
		jwtExp: jwtExp,
	}
}

// jwtClaims extends standard JWT claims with application-specific user ID.
type jwtClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"id"`
}

// Validate implements custom claims validation logic.
func (c jwtClaims) Validate() error {
	if c.UserID == "" {
		return errNoUserID
	}
	return nil
}

// NewJWTString creates a new JWT token for the given user ID.
//
// The token includes standard claims (exp) and stores userID in sub claim.
// Returns the signed token string or error if signing fails.
func (s *Auth) NewJWTString(userID models.UserID) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExp)),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseIDFromAuthHeader checks JWT token validity and returns contained user ID.
//
// Verifies token signature and expiration time. Returns userID if token is valid
// or error describing validation failure.
func (s *Auth) ParseIDFromAuthHeader(header string) (models.UserID, error) {
	tokenString := strings.TrimPrefix(header, "Bearer")
	tokenString = strings.TrimSpace(tokenString)

	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// GetUserIDFromCtx extracts user ID from context set by Auth middleware.
//
// Typical usage in HTTP handlers:
//
//	userID, err := auth.GetUserIDFromCtx(r.Context())
//
// Returns empty string and error if user is not authenticated.
func (s *Auth) GetUserIDFromCtx(ctx context.Context) (models.UserID, error) {
	uid, ok := ctx.Value(contextkeys.UserID).(models.UserID)
	if !ok {
		return "", errNoUserID
	}
	return uid, nil
}
