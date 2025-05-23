package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKey = "secret"

var testJWTExp = time.Duration(1) * time.Minute

func TestAuth_NewJWTString(t *testing.T) {
	jwtService := NewAuth(testKey, testJWTExp)

	t.Run("valid test", func(t *testing.T) {
		jwtString, err := jwtService.NewJWTString(testUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, jwtString)
	})
}

func TestParseIDFromAuthHeader(t *testing.T) {
	jwtService := NewAuth(testKey, testJWTExp)

	t.Run("valid test", func(t *testing.T) {
		claims := jwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(testJWTExp)),
			},
			UserID: testUserID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testKey))
		require.NoError(t, err)
		authHeader := "Bearer " + tokenString
		uid, err := jwtService.ParseIDFromAuthHeader(authHeader)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("no user id", func(t *testing.T) {
		claims := jwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(testJWTExp)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testKey))
		authHeader := "Bearer " + tokenString
		require.NoError(t, err)
		_, err = jwtService.ParseIDFromAuthHeader(authHeader)
		assert.Error(t, err)
	})
}

func TestAuth_GetUserIDFromCtx(t *testing.T) {
	jwtService := NewAuth(testKey, testJWTExp)

	t.Run("valid test", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextkeys.UserID, testUserID)
		uid, err := jwtService.GetUserIDFromCtx(ctx)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("no user id error", func(t *testing.T) {
		_, err := jwtService.GetUserIDFromCtx(context.Background())
		assert.ErrorIs(t, err, errNoUserID)
	})
}
