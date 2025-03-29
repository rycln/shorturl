package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExp = time.Hour * 24

var ErrNoToken = errors.New("no jwt token")

type jwtClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

func getTokenAndUID(c *fiber.Ctx, key string) (string, string, error) {
	rawToken := string(c.Request().Header.Peek("Authorization"))
	if rawToken == "" {
		return "", "", ErrNoToken
	}
	rawToken = strings.TrimPrefix(rawToken, "Bearer")
	rawToken = strings.TrimSpace(rawToken)
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", "", err
	}
	return rawToken, claims.ID, nil
}

func makeUserID() string {
	return uuid.NewString()
}

func makeTokenString(uid, key string) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		ID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
