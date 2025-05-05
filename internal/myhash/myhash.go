package myhash

import (
	"crypto/md5"

	"github.com/jxskiss/base62"
)

const (
	shortURLSize = 7
)

func Base62(url string) string {
	hash := md5.Sum([]byte(url))
	encodedHash := base62.Encode([]byte(hash[:]))
	return string(encodedHash[:shortURLSize])
}
