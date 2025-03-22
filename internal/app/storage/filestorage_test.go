package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorageAddURL(t *testing.T) {
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
			fileName := "test"
			strg, err := NewFileStorage(fileName)
			require.NoError(t, err)
			defer os.Remove(fileName)
			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				var err error
				for i := range test.shortURLs {
					surl := NewShortenedURL(testID, test.shortURLs[i], test.origURLs[i])
					err = strg.AddURL(context.Background(), surl)
				}
				if test.want.wantErr {
					assert.Error(t, err)
				}
				for k, v := range test.want.mustContain {
					var origURL string
					fd, err := newFileDecoder(fileName)
					require.NoError(t, err)
					defer fd.close()
					for {
						surl := &ShortenedURL{}
						err := fd.decoder.Decode(surl)
						require.NoError(t, err)
						if surl.ShortURL == k {
							origURL = surl.OrigURL
							break
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

func TestFileStorageAddBatchURL(t *testing.T) {
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
			fileName := "test"
			strg, err := NewFileStorage(fileName)
			require.NoError(t, err)
			defer os.Remove(fileName)
			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				var surls = make([]ShortenedURL, len(test.shortURLs))
				var err error
				for i := range test.shortURLs {
					surl := NewShortenedURL(testID, test.shortURLs[i], test.origURLs[i])
					surls[i] = surl
				}
				err = strg.AddBatchURL(context.Background(), surls)
				if test.want.wantErr {
					assert.Error(t, err)
				}
				for k, v := range test.want.mustContain {
					var origURL string
					fd, err := newFileDecoder(fileName)
					require.NoError(t, err)
					defer fd.close()
					for {
						surl := &ShortenedURL{}
						err := fd.decoder.Decode(surl)
						require.NoError(t, err)
						if surl.ShortURL == k {
							origURL = surl.OrigURL
							break
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

func TestFileStorageGetOrigURL(t *testing.T) {
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileName := "test"
			strg, err := NewFileStorage(fileName)
			require.NoError(t, err)
			defer os.Remove(fileName)

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					surl := NewShortenedURL(testID, test.shortURLs[i], test.origURLs[i])
					err := strg.encoder.encoder.Encode(&surl)
					require.NoError(t, err)
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

func TestFileStorageGetShortURL(t *testing.T) {
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileName := "test"
			strg, err := NewFileStorage(fileName)
			require.NoError(t, err)
			defer os.Remove(fileName)

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					surl := NewShortenedURL(testID, test.shortURLs[i], test.origURLs[i])
					err := strg.encoder.encoder.Encode(&surl)
					require.NoError(t, err)
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

func TestFileStorageGetAllUserURLs(t *testing.T) {
	type want struct {
		mustContain map[string]string
		uid         string
		wantErr     bool
	}

	tests := []struct {
		name      string
		dataUID   string
		shortURLs []string
		origURLs  []string
		want      want
	}{
		{
			name:    "Simple test #1",
			dataUID: testID,
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
				uid:     testID,
				wantErr: false,
			},
		},
		{
			name:    "Simple test #2",
			dataUID: testID,
			shortURLs: []string{
				"1234ABC",
				"235DCE",
			},
			origURLs: []string{
				"https://ya.ru/",
				"https://yandex.ru/",
			},
			want: want{
				mustContain: map[string]string{
					"1234ABC": "https://ya.ru/",
					"235DCE":  "https://yandex.ru/",
				},
				uid:     testID,
				wantErr: false,
			},
		},
		{
			name:    "Wrong GetURL request #1",
			dataUID: "testID",
			shortURLs: []string{
				"abcdefg",
			},
			origURLs: []string{
				"https://practicum.yandex.ru/",
			},
			want: want{
				uid:     "2",
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileName := "test"
			strg, err := NewFileStorage(fileName)
			require.NoError(t, err)
			defer os.Remove(fileName)

			if assert.Equal(t, len(test.shortURLs), len(test.origURLs), "wrong tests") {
				for i := range test.shortURLs {
					surl := NewShortenedURL(test.dataUID, test.shortURLs[i], test.origURLs[i])
					err := strg.encoder.encoder.Encode(&surl)
					require.NoError(t, err)
				}
				surls, err := strg.GetAllUserURLs(context.Background(), test.want.uid)
				if err != nil {
					if test.want.wantErr {
						assert.Error(t, err)
					}
				}
				if !test.want.wantErr {
					for _, surl := range surls {
						assert.Equal(t, test.want.mustContain[surl.ShortURL], surl.OrigURL)
					}
				}
			}
		})
	}
}
