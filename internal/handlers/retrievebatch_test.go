package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/rycln/shorturl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveBatch_Handle(t *testing.T) {
	type want struct {
		code        int
		resContains string
		wantErr     bool
	}
	tests := []struct {
		name          string
		method        string
		path          string
		dataUID       string
		reqUID        string
		storeContains []storage.ShortenedURL
		want          want
	}{
		{
			name:    "Valid test #1",
			method:  http.MethodGet,
			path:    "/api/user/urls",
			dataUID: testID,
			reqUID:  testID,
			storeContains: []storage.ShortenedURL{
				{
					UserID:   testID,
					ShortURL: "abcd",
					OrigURL:  "https://practicum.yandex.ru/",
				},
			},
			want: want{
				code:        http.StatusOK,
				resContains: `"short_url":"http://localhost:8080/abcd"`,
				wantErr:     false,
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodPost,
			path:   "/api/user/urls",
			want: want{
				code: http.StatusBadRequest,
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
		{
			name:    "Request without token #1",
			method:  http.MethodGet,
			path:    "/api/user/urls",
			dataUID: testID,
			reqUID:  testID,
			storeContains: []storage.ShortenedURL{
				{
					UserID:   testID,
					ShortURL: "abcd",
					OrigURL:  "https://practicum.yandex.ru/",
				},
			},
			want: want{
				code:    http.StatusNoContent,
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mocks.NewMockshortenConfiger(ctrl)
			mStrg := mocks.NewMockretrieveBatchStorager(ctrl)

			if test.want.code != http.StatusBadRequest {
				mCfg.EXPECT().GetKey().Return(testKey)
				if !test.want.wantErr {
					mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
					mStrg.EXPECT().GetAllUserURLs(gomock.Any(), test.dataUID).Return(test.storeContains, nil)
				}
			}

			app := fiber.New()
			app.Get("/api/user/urls", timeout.NewWithContext(NewRetrieveBatchHandler(mStrg, mCfg), testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			request := httptest.NewRequest(test.method, test.path, nil)
			if !test.want.wantErr {
				token, err := makeTokenString(test.reqUID, testKey)
				require.NoError(t, err)
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			require.Equal(t, test.want.code, res.StatusCode)
			if test.want.code != http.StatusBadRequest {
				if !test.want.wantErr {
					resBody, err := io.ReadAll(res.Body)
					require.NoError(t, err)
					assert.Contains(t, string(resBody), test.want.resContains)
				} else {
					require.NotEmpty(t, res.Cookies())
				}
			}
		})
	}
}
