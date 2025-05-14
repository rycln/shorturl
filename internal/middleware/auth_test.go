package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/middleware/mocks"
	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = models.UserID("1")
	testJWT    = "abc.123.def"
)

var errTest = errors.New("test error")

func TestAuthMiddleware_JWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mJWT := mocks.NewMockauthServicer(ctrl)

	auth := NewAuthMiddleware(mJWT)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value(contextkeys.UserID).(models.UserID)
		require.True(t, ok)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(uid))
	})

	t.Run("valid test", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(gomock.Any()).Return(testUserID, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Authorization", testJWT)
		w := httptest.NewRecorder()

		auth.JWT(testHandler).ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(testUserID), string(resBody))
	})

	t.Run("no jwt", func(t *testing.T) {
		mJWT.EXPECT().NewJWTString(gomock.Any()).Return(testJWT, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		auth.JWT(testHandler).ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.NotEmpty(t, resBody)
		assert.NotEmpty(t, res.Cookies())
	})

	t.Run("invalid jwt", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(gomock.Any()).Return(models.UserID(""), errTest)
		mJWT.EXPECT().NewJWTString(gomock.Any()).Return(testJWT, nil)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Add("Authorization", testJWT)
		w := httptest.NewRecorder()

		auth.JWT(testHandler).ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.NotEmpty(t, resBody)
		assert.NotEmpty(t, res.Cookies())
	})

	t.Run("new jwt error", func(t *testing.T) {
		mJWT.EXPECT().NewJWTString(gomock.Any()).Return("", errTest)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		auth.JWT(testHandler).ServeHTTP(w, request)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
