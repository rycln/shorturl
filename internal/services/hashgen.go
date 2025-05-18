package services

import (
	"crypto/md5"

	"github.com/jxskiss/base62"
	"github.com/rycln/shorturl/internal/models"
)

// HashGen generates unique hash strings from original URLs.
//
// The service creates consistent, short hash representations suitable for URL shortening.
type HashGen struct {
	len int
}

// NewHashGen creates a new HashGen instance.
func NewHashGen(len int) *HashGen {
	return &HashGen{
		len: len,
	}
}

// GenerateHashFromURL creates a hash string from the original URL.
//
// The same URL will always produce the same hash.
func (s *HashGen) GenerateHashFromURL(orig models.OrigURL) models.ShortURL {
	hash := md5.Sum([]byte(orig))
	encodedHash := base62.Encode([]byte(hash[:]))
	return models.ShortURL(encodedHash[:s.len])
}
