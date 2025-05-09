package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testHashLen = 10

func TestHashGen_GenerateHashFromURL(t *testing.T) {
	s := NewHashGen(testHashLen)

	t.Run("valid test", func(t *testing.T) {
		hash := s.GenerateHashFromURL(testOrigURL)
		assert.Len(t, hash, testHashLen)
	})
}
