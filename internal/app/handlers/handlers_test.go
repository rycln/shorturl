package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerVariables_ShortenURL(t *testing.T) {
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			store := mem.NewSimpleMemStorage()
			hv := NewHandlerVariables(store)
			app.All("/", hv.ShortenURL)

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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			store := mem.NewSimpleMemStorage()
			hv := NewHandlerVariables(store)
			app.All("/:short", hv.ReturnURL)
			for shortURL, fullURL := range test.storeContains {
				hv.store.AddURL(shortURL, fullURL)
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
