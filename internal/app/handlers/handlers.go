package handlers

import (
	"strings"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func getUserID(c *fiber.Ctx) (string, error) {
	rawToken := string(c.Request().Header.Peek("Authorization"))
	if rawToken == "" {
		return makeUserID(), nil
	}
	rawToken = strings.TrimPrefix(rawToken, "Bearer")
	rawToken = strings.TrimSpace(rawToken)
	token, err := jwt.ParseWithClaims(rawToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
}

func makeUserID() string {
	return uuid.NewString()
}
