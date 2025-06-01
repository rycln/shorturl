package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/handlers/mocks"
	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mShort := mocks.NewMockshortenServicer(ctrl)
	mAuth := mocks.NewMockshortenAuthServicer(ctrl)

	shortenHandler := NewShortenHandler(mShort, mAuth, testBaseAddr)

	t.Run("valid test", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mShort.EXPECT().ShortenURL(gomock.Any(), testPair.UID, testPair.Orig).Return(&testPair, nil)

		reqBody := strings.NewReader(string(testPair.Orig))
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		wantBody := testBaseAddr + "/" + string(testPair.Short)
		assert.Equal(t, wantBody, string(resBody))

	})

	t.Run("user id error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(models.UserID(""), errTest)

		reqBody := strings.NewReader(string(testPair.Orig))
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("wrong orig url error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)

		reqBody := strings.NewReader("wrong url")
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("conflict", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mErr := mocks.NewMockerrShortenConflict(ctrl)
		mErr.EXPECT().IsErrConflict().Return(true)
		mShort.EXPECT().ShortenURL(gomock.Any(), testPair.UID, testPair.Orig).Return(&testPair, mErr)

		reqBody := strings.NewReader(string(testPair.Orig))
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusConflict, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		wantBody := testBaseAddr + "/" + string(testPair.Short)
		assert.Equal(t, wantBody, string(resBody))
	})

	t.Run("some error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mShort.EXPECT().ShortenURL(gomock.Any(), testPair.UID, testPair.Orig).Return(nil, errTest)

		reqBody := strings.NewReader(string(testPair.Orig))
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
