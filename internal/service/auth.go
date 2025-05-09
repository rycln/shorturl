package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rycln/shorturl/internal/models"
)

var ErrNoUserID = errors.New("jwt does not contain user id")

type AuthService struct {
	jwtKey string
	jwtExp time.Duration
}

func NewAuthService(jwtkey string, jwtExp time.Duration) *AuthService {
	return &AuthService{
		jwtKey: jwtkey,
		jwtExp: jwtExp,
	}
}

func (s *AuthService) MakeUserID() string {
	return uuid.NewString()
}

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"id"`
}

func (c jwtClaims) Validate() error {
	if c.UserID == "" {
		return ErrNoUserID
	}
	return nil
}

func (s *AuthService) NewJWTString(userID models.UserID) (string, error) {
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

func (s *AuthService) ParseIDFromAuthHeader(header string) (models.UserID, error) {
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
