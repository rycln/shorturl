package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddURL(t *testing.T) {
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
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strg := NewSimpleStorage()
			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				var err error
				for i := range test.shortURLs {
					surl := NewShortenedURL(test.shortURLs[i], test.origURLs[i])
					err = strg.AddURL(context.Background(), surl)
				}
				if test.want.wantErr {
					assert.Error(t, err)
				}
				for k, v := range test.want.mustContain {
					origURL := strg.storage[k]
					if !test.want.wantErr {
						assert.Equal(t, v, origURL)
					}
				}
			}
		})
	}
}

func TestAddBatchURL(t *testing.T) {
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strg := NewSimpleStorage()
			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				var surls = make([]ShortenedURL, len(test.shortURLs))
				var err error
				for i := range test.shortURLs {
					surl := NewShortenedURL(test.shortURLs[i], test.origURLs[i])
					surls[i] = surl
				}
				err = strg.AddBatchURL(context.Background(), surls)
				if test.want.wantErr {
					assert.Error(t, err)
				}
				for k, v := range test.want.mustContain {
					origURL := strg.storage[k]
					if !test.want.wantErr {
						assert.Equal(t, v, origURL)
					}
				}
			}
		})
	}
}

func TestGetOrigURL(t *testing.T) {
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
			name: "Wrong GetURL request #1",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strg := NewSimpleStorage()

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					strg.storage[test.shortURLs[i]] = test.origURLs[i]
				}
				for k, v := range test.want.mustContain {
					origURL, err := strg.GetOrigURL(context.Background(), k)
					if err != nil {
						if test.want.wantErr {
							assert.Error(t, err)
						}
					}
					if !test.want.wantErr {
						assert.Equal(t, v, origURL)
					}
				}
			}
		})
	}
}

func TestShortOrigURL(t *testing.T) {
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
			name: "Wrong GetURL request #1",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strg := NewSimpleStorage()

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					strg.storage[test.shortURLs[i]] = test.origURLs[i]
				}
				for k, v := range test.want.mustContain {
					shortURL, err := strg.GetShortURL(context.Background(), v)
					if err != nil {
						if test.want.wantErr {
							assert.Error(t, err)
						}
					}
					if !test.want.wantErr {
						assert.Equal(t, k, shortURL)
					}
				}
			}
		})
	}
}
