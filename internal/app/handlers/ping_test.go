package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/stretchr/testify/require"
)

func TestPing_Handle(t *testing.T) {
	type want struct {
		code    int
		wantErr bool
	}
	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			name:   "Valid test #1",
			method: http.MethodGet,
			path:   "/ping",
			want: want{
				code:    http.StatusOK,
				wantErr: false,
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodPost,
			path:   "/ping",
			want: want{
				code:    http.StatusBadRequest,
				wantErr: false,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodPost,
			path:   "/pin",
			want: want{
				code:    http.StatusBadRequest,
				wantErr: false,
			},
		},
		{
			name:   "Ping failed #1",
			method: http.MethodGet,
			path:   "/ping",
			want: want{
				code:    http.StatusInternalServerError,
				wantErr: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mStrg := mocks.NewMockpingStorager(ctrl)

			if test.want.code != http.StatusBadRequest {
				if !test.want.wantErr {
					mStrg.EXPECT().Ping(gomock.Any()).Return(nil)
				} else {
					mStrg.EXPECT().Ping(gomock.Any()).Return(errTest)
				}
			}

			app := fiber.New()
			app.Get("/ping", timeout.NewWithContext(NewPingHandler(mStrg), testTimeoutDuration))
			app.Use(func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusBadRequest)
			})

			request := httptest.NewRequest(test.method, test.path, nil)

			res, err := app.Test(request, -1)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			require.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
