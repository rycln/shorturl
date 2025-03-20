package handlers

import (
	"strings"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

func getUserID(c *fiber.Ctx, key string) string {
	rawToken := string(c.Request().Header.Peek("Authorization"))
	if rawToken == "" {
		return makeUserID()
	}
	rawToken = strings.TrimPrefix(rawToken, "Bearer")
	rawToken = strings.TrimSpace(rawToken)
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return makeUserID()
	}
	return claims.ID
}

func makeUserID() string {
	return uuid.NewString()
}
