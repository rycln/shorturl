package service

import (
	"crypto/md5"

	"github.com/jxskiss/base62"
	"github.com/rycln/shorturl/internal/models"
)

type HashGen struct {
	len int
}

func NewHashGen(len int) *HashGen {
	return &HashGen{
		len: len,
	}
}

func (s *HashGen) GenerateHashFromURL(orig models.OrigURL) models.ShortURL {
	hash := md5.Sum([]byte(orig))
	encodedHash := base62.Encode([]byte(hash[:]))
	return models.ShortURL(encodedHash[:s.len])
}
