package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStorager interface {
	AddURL(string, string)
	GetURL(string) (string, error)
}

type testStorage struct {
	ts testStorager
}

func NewTestStorage(ts testStorager) *testStorage {
	return &testStorage{
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
		fullURLs  []string
		want      want
	}{
		{
			name: "Simple test #1",
			shortURLs: []string{
				"abcdefg",
			},
			fullURLs: []string{
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
			fullURLs: []string{
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
			fullURLs: []string{
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
			fullURLs: []string{
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
			fullURLs: []string{
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
			storage := NewSimpleMemStorage()
			ts := NewTestStorage(storage)

			if assert.Equal(t, len(test.shortURLs), len(test.fullURLs), "wrong tests") {
				for i := range test.shortURLs {
					ts.ts.AddURL(test.shortURLs[i], test.fullURLs[i])
				}
				for k, v := range test.want.mustContain {
					recFullURL, err := ts.ts.GetURL(k)
					if err != nil {
						if !test.want.wantErr {
							assert.Error(t, err)
						}
					}
					if !test.want.wantErr {
						assert.Equal(t, v, recFullURL)
					}
				}
			}
		})
	}
}
