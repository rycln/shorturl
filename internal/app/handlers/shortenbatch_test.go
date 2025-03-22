package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenBatch_Handle(t *testing.T) {
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
			path:   "/api/shorten/batch",
			body:   []byte(`[ {"correlation_id":"abc","original_url":"https://practicum.yandex.ru/"} ]`),
			want: want{
				code:        http.StatusCreated,
				resContains: `"correlation_id":"abc"`,
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodGet,
			path:   "/api/shorten/batch",
			body:   []byte(`[ {"correlation_id":"abc","original_url":"https://practicum.yandex.ru/"} ]`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong URL #1",
			method: http.MethodPost,
			path:   "/api/shorten/batc",
			body:   []byte(`[ {"correlation_id":"abc","original_url":"https://practicum.yandex.ru/"} ]`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong JSON #1",
			method: http.MethodPost,
			path:   "/api/shorten/batch",
			body:   []byte(`[ {"correlation_id":"abc","original_url":"https://practicum.yandex.ru/"`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodPost,
			path:   "/api/shorten/bad",
			body:   []byte(`[ {"correlation_id":"abc","original_url":"https://practicum.yandex.ru/"} ]`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mocks.NewMockbatchConfiger(ctrl)
			mStrg := mocks.NewMockbatchStorager(ctrl)

			if test.want.code != http.StatusBadRequest {
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mCfg.EXPECT().GetKey().Return(testKey)
				mStrg.EXPECT().AddBatchURL(gomock.Any(), gomock.Any()).Return(nil)
			}

			app := fiber.New()
			app.Post("/api/shorten/batch", timeout.NewWithContext(NewShortenBatch(mStrg, mCfg, testHash).Handle, testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			bodyReader := bytes.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.path, bodyReader)
			request.Header.Set("Content-Type", test.want.contentType)

			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			require.Equal(t, test.want.code, res.StatusCode)
			if test.want.code != http.StatusBadRequest {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Contains(t, string(resBody), test.want.resContains)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
