package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/app/mocks"
	"github.com/rycln/shorturl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBaseAddr        = "http://localhost:8080"
	testHashVal         = "abc"
	testTimeoutDuration = time.Duration(2) * time.Second
)

var errTest = errors.New("test error")

func testHash(str string) string {
	return testHashVal
}

func TestServerArgs_ShortenURL(t *testing.T) {
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

			mCfg := mocks.NewMockconfiger(ctrl)
			mStrg := mocks.NewMockstorager(ctrl)

			mCfg.EXPECT().TimeoutDuration().Return(testTimeoutDuration).AnyTimes()
			switch test.want.code {
			case http.StatusCreated:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(nil)
			case http.StatusConflict:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(storage.ErrConflict)
				mStrg.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return(testHashVal, nil)
			}

			app := fiber.New()
			app.Post("/", timeout.NewWithContext(shorten.Handle, to))
			sa := NewServerArgs(mStrg, mCfg, testHash)
			Set(app, sa)

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

func TestServerArgs_ReturnURL(t *testing.T) {
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

			mCfg := mocks.NewMockconfiger(ctrl)
			mStrg := mocks.NewMockstorager(ctrl)

			mCfg.EXPECT().TimeoutDuration().Return(testTimeoutDuration).AnyTimes()
			testShortURL := strings.TrimPrefix(test.path, "/")
			if test.want.code != http.StatusBadRequest {
				mStrg.EXPECT().GetOrigURL(gomock.Any(), testShortURL).Return(test.storeContains[testShortURL], nil)
			}
			if test.want.wantErr {
				mStrg.EXPECT().GetOrigURL(gomock.Any(), testShortURL).Return(test.storeContains[testShortURL], errTest)
			}

			app := fiber.New()
			sa := NewServerArgs(mStrg, mCfg, testHash)
			Set(app, sa)

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

func TestServerArgs_ShortenAPI(t *testing.T) {
	type want struct {
		code        int
		resEqual    string
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
				resEqual:    fmt.Sprintf(`{"result":"%s/%s"}`, testBaseAddr, testHashVal),
				contentType: "application/json",
			},
		},
		{
			name:   "Wrong method #1",
			method: http.MethodGet,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong URL #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"practicum"}`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong content type #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong JSON #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url:"https://practicum.yandex.ru/"}`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Wrong path #1",
			method: http.MethodPost,
			path:   "/api/shorten/bad",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "Same URL sended twice #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   []byte(`{"url":"https://practicum.yandex.ru/"}`),
			want: want{
				code:        http.StatusConflict,
				resEqual:    fmt.Sprintf(`{"result":"%s/%s"}`, testBaseAddr, testHashVal),
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mCfg := mocks.NewMockconfiger(ctrl)
			mStrg := mocks.NewMockstorager(ctrl)

			mCfg.EXPECT().TimeoutDuration().Return(testTimeoutDuration).AnyTimes()
			switch test.want.code {
			case http.StatusCreated:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(nil)
			case http.StatusConflict:
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mStrg.EXPECT().AddURL(gomock.Any(), gomock.Any()).Return(storage.ErrConflict)
				mStrg.EXPECT().GetShortURL(gomock.Any(), gomock.Any()).Return(testHashVal, nil)
			}

			app := fiber.New()
			sa := NewServerArgs(mStrg, mCfg, testHash)
			Set(app, sa)

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
				assert.Contains(t, string(resBody), test.want.resEqual)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestServerArgs_ShortenBatch(t *testing.T) {
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

			mCfg := mocks.NewMockconfiger(ctrl)
			mStrg := mocks.NewMockstorager(ctrl)

			mCfg.EXPECT().TimeoutDuration().Return(testTimeoutDuration).AnyTimes()
			if test.want.code != http.StatusBadRequest {
				mCfg.EXPECT().GetBaseAddr().Return(testBaseAddr)
				mStrg.EXPECT().AddBatchURL(gomock.Any(), gomock.Any()).Return(nil)
			}

			app := fiber.New()
			sa := NewServerArgs(mStrg, mCfg, testHash)
			Set(app, sa)

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

func TestServerArgs_PingDB(t *testing.T) {
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

			mCfg := mocks.NewMockconfiger(ctrl)
			mStrg := mocks.NewMockstorager(ctrl)

			mCfg.EXPECT().TimeoutDuration().Return(testTimeoutDuration).AnyTimes()
			if test.want.code != http.StatusBadRequest {
				if !test.want.wantErr {
					mStrg.EXPECT().Ping(gomock.Any()).Return(nil)
				} else {
					mStrg.EXPECT().Ping(gomock.Any()).Return(errTest)
				}
			}

			app := fiber.New()
			sa := NewServerArgs(mStrg, mCfg, testHash)
			Set(app, sa)

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
