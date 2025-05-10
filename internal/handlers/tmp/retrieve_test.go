package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieve_Handle(t *testing.T) {
	type want struct {
		code     int
		location string
		wantErr  bool
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
				code:    http.StatusBadRequest,
				wantErr: true,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodGet,
			path:   "/dcba/abcd",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mStrg := mocks.NewMockretrieveStorager(ctrl)

			testShortURL := strings.TrimPrefix(test.path, "/")
			if test.want.code != http.StatusBadRequest {
				mStrg.EXPECT().GetOrigURL(gomock.Any(), testShortURL).Return(test.storeContains[testShortURL], nil)
			}
			if test.want.wantErr {
				mStrg.EXPECT().GetOrigURL(gomock.Any(), testShortURL).Return(test.storeContains[testShortURL], errTest)
			}

			app := fiber.New()
			app.Get("/:short", timeout.NewWithContext(NewRetrieveHandler(mStrg), testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			request := httptest.NewRequest(test.method, test.path, nil)
			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			res.Body.Close()
			require.Equal(t, test.want.code, res.StatusCode)
			if test.want.code != http.StatusBadRequest {
				assert.Equal(t, test.want.location, res.Header.Get("Location"))
			}
		})
	}
}
