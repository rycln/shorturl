package myhash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const base62Dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func TestBase62(t *testing.T) {
	tests := []struct {
		name string
		url  string
		iter int
	}{
		{
			name: "Simple test #1",
			url:  "https://practicum.yandex.ru/",
		},
		{
			name: "Simple test #2",
			url:  "https://ya.ru/",
		},
		{
			name: "Same URL several times #1",
			url:  "https://practicum.yandex.ru/",
			iter: 5,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encVal := Base62(test.url)
			for i := 0; i < test.iter; i++ {
				assert.Equal(t, encVal, Base62(test.url))
			}
			assert.Len(t, encVal, shortURLSize)
			assert.Subset(t, []byte(base62Dict), []byte(encVal))
		})
	}
}
