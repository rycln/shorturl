package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/stretchr/testify/require"
)

func TestDeleteBatch_Handle(t *testing.T) {
	type want struct {
		code    int
		errWant bool
	}
	tests := []struct {
		name    string
		method  string
		path    string
		reqBody []byte
		want    want
	}{
		{
			name:    "Valid test #1",
			method:  http.MethodDelete,
			path:    "/api/user/urls",
			reqBody: []byte(`[ "abc123", "def456" ]`),
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodPatch,
			path:   "/api/user/urls",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodDelete,
			path:   "/dcba/abcd",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "Wrong JSON #1",
			method:  http.MethodDelete,
			path:    "/api/user/urls",
			reqBody: []byte(`[ "abc123", "def456" `),
			want: want{
				code:    http.StatusBadRequest,
				errWant: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mocks.NewMockdeleteBatchConfiger(ctrl)
			mStrg := mocks.NewMockdeleteBatchStorager(ctrl)

			if test.want.code != http.StatusBadRequest || test.want.errWant {
				mCfg.EXPECT().GetKey().Return(testKey)
			}

			app := fiber.New()
			app.Delete("/api/user/urls", timeout.NewWithContext(NewDeleteBatchHandler(context.Background(), mStrg, mCfg), testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			bodyReader := bytes.NewReader(test.reqBody)
			request := httptest.NewRequest(test.method, test.path, bodyReader)
			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			require.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
