package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/shorturl/internal/models"
)

const tokenExp = time.Duration(2) * time.Hour

var ErrNoUserID = errors.New("jwt does not contain user id")

type JWTService struct {
	key string
}

func NewJWTService(key string) *JWTService {
	return &JWTService{
		key: key,
	}
}

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"id"`
}

func (c jwtClaims) Validate() error {
	if c.UserID == 0 {
		return ErrNoUserID
	}
	return nil
}

func (s *JWTService) NewJWTString(userID models.UserID) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) ParseIDFromAuthHeader(header string) (models.UserID, error) {
	tokenString := strings.TrimPrefix(header, "Bearer")
	tokenString = strings.TrimSpace(tokenString)

	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
