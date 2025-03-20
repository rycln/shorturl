package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/rycln/shorturl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShorten_Handle(t *testing.T) {
	type want struct {
		code        int
		resEqual    string
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
				resEqual:    testBaseAddr + "/" + testHashVal,
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
		{
			name:   "Same URL sended twice #1",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusConflict,
				resEqual:    testBaseAddr + "/" + testHashVal,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mocks.NewMockshortenConfiger(ctrl)
			mStrg := mocks.NewMockshortenStorager(ctrl)

			switch test.want.code {
			case http.StatusCreated:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mCfg.EXPECT().GetKey().Return(testKey)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(nil)
			case http.StatusConflict:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mCfg.EXPECT().GetKey().Return(testKey)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(storage.ErrConflict)
				mStrg.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return(testHashVal, nil)
			}

			app := fiber.New()
			app.Post("/", timeout.NewWithContext(NewShorten(mStrg, mCfg, testHash).Handle, testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			bodyReader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.path, bodyReader)

			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			require.Equal(t, test.want.code, res.StatusCode)
			if test.want.code != http.StatusBadRequest {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, string(resBody), test.want.resEqual)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
