package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStorager interface {
	AddURL(context.Context, ...ShortenedURL) error
	GetURL(context.Context, string) (string, error)
}

type testStorage struct {
	ts testStorager
}

func NewTestStorage(ts testStorager) testStorage {
	return testStorage{
		ts: ts,
	}
}

func TestAddURLAndGetURL(t *testing.T) {
	type want struct {
		mustContain map[string]string
		wantErr     bool
	}

	tests := []struct {
		name      string
		shortURLs []string
		origURLs  []string
		want      want
	}{
		{
			name: "Simple test #1",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"abcdefg": "https://practicum.yandex.ru/",
				},
				wantErr: false,
			},
		},
		{
			name: "Simple test #2",
			shortURLs: []string{
				"1234ABC",
			},
			origURLs: []string{
				"https://ya.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"1234ABC": "https://ya.ru/",
				},
				wantErr: false,
			},
		},
		{
			name: "Several pairs of data #1",
			shortURLs: []string{
				"abcdefg",
				"1234ABC",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
				"https://ya.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"abcdefg": "https://practicum.yandex.ru/",
					"1234ABC": "https://ya.ru/",
				},
				wantErr: false,
			},
		},
		{
			name: "Same data #1",
			shortURLs: []string{
				"abcdefg",
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
				"https://practicum.yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"abcdefg": "https://practicum.yandex.ru/",
				},
				wantErr: false,
			},
		},
		{
			name: "Wrong GetURL request #1",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"1234ABC": "https://ya.ru/",
				},
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strg := NewSimpleStorage()
			ts := NewTestStorage(strg)

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					surl := NewShortenedURL(test.shortURLs[i], test.origURLs[i])
					ts.ts.AddURL(context.Background(), surl)
				}
				for k, v := range test.want.mustContain {
					recOrigURL, err := ts.ts.GetURL(context.Background(), k)
					if err != nil {
						if !test.want.wantErr {
							assert.Error(t, err)
						}
					}
					if !test.want.wantErr {
						assert.Equal(t, v, recOrigURL)
					}
				}
			}
		})
	}
}
