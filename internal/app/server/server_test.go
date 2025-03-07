package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerVariables_ShortenURL(t *testing.T) {
	config := &config.Cfg{
		ServerAddr:    config.DefaultServerAddr,
		ShortBaseAddr: config.DefaultBaseAddr,
	}

	app := fiber.New()
	strg := storage.NewSimpleStorage()
	sa := NewServerArgs(strg, config)
	Set(app, sa)

	type want struct {
		code        int
		resContains string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			name:   "Valid test #1",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusCreated,
				resContains: "http://localhost:8080/",
				contentType: "text/plain",
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodGet,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong URL #1",
			method: http.MethodPost,
			path:   "/",
			body:   "practicum",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodPost,
			path:   "/ab/cd",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.path, bodyReader)

			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			require.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode != http.StatusBadRequest {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Contains(t, string(resBody), test.want.resContains)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestHandlerVariables_ReturnURL(t *testing.T) {
	config := &config.Cfg{
		ServerAddr:    config.DefaultServerAddr,
		ShortBaseAddr: config.DefaultBaseAddr,
	}

	app := fiber.New()
	strg := storage.NewSimpleStorage()
	sa := NewServerArgs(strg, config)
	Set(app, sa)

	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name          string
		method        string
		path          string
		storeContains map[string]string
		want          want
	}{
		{
			name:   "Valid test #1",
			method: http.MethodGet,
			path:   "/abcd",
			storeContains: map[string]string{
				"abcd": "https://practicum.yandex.ru/",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://practicum.yandex.ru/",
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodPost,
			path:   "/abcd",
			storeContains: map[string]string{
				"abcd": "https://practicum.yandex.ru/",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Short URL does not exist #1",
			method: http.MethodGet,
			path:   "/dcba",
			storeContains: map[string]string{
				"abcd": "https://practicum.yandex.ru/",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodGet,
			path:   "/dcba/abcd",
			storeContains: map[string]string{
				"abcd": "https://practicum.yandex.ru/",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for shortURL, origURL := range test.storeContains {
				surl := storage.NewShortenedURL(shortURL, origURL)
				sa.strg.AddURL(context.Background(), surl)
			}
			request := httptest.NewRequest(test.method, test.path, nil)
			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			res.Body.Close()
			require.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.location, res.Header.Get("Location"))
			}
		})
	}
}

func TestServerArgs_ShortenAPI(t *testing.T) {
	config := &config.Cfg{
		ServerAddr:    config.DefaultServerAddr,
		ShortBaseAddr: config.DefaultBaseAddr,
	}

	app := fiber.New()
	strg := storage.NewSimpleStorage()
	sa := NewServerArgs(strg, config)
	Set(app, sa)

	type want struct {
		code        int
		resContains string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
		want   want
	}{
		{
			name:   "Valid test #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusCreated,
				resContains: "http://localhost:8080/",
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodGet,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong URL #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"practicum"}`),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong content type #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:   "Wrong JSON #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url:"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodPost,
			path:   "/api/shorten/bad",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.path, bodyReader)
			request.Header.Set("Content-Type", test.want.contentType)

			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			require.Equal(t, test.want.code, res.StatusCode)
			if res.StatusCode != http.StatusBadRequest {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Contains(t, string(resBody), test.want.resContains)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
