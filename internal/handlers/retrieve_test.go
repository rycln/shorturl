package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/handlers/mocks"
	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mServ := mocks.NewMockretrieveServicer(ctrl)

	retrieveHandler := NewRetrieveHandler(mServ)

	t.Run("valid test", func(t *testing.T) {
		mServ.EXPECT().GetShortURLFromCtx(gomock.Any()).Return(testShortURL, nil)
		mServ.EXPECT().GetOrigURLByShort(gomock.Any(), testShortURL).Return(testOrigURL, nil)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, res.Header.Get("Location"), string(testOrigURL))
	})

	t.Run("short url error", func(t *testing.T) {
		mServ.EXPECT().GetShortURLFromCtx(gomock.Any()).Return(models.ShortURL(""), errTest)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("url was deleted", func(t *testing.T) {
		mServ.EXPECT().GetShortURLFromCtx(gomock.Any()).Return(testShortURL, nil)
		mErr := mocks.NewMockerrRetrieveDeletedURL(ctrl)
		mErr.EXPECT().IsErrDeletedURL().Return(true)
		mServ.EXPECT().GetOrigURLByShort(gomock.Any(), testShortURL).Return(models.OrigURL(""), mErr)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusGone, res.StatusCode)
	})

	t.Run("some service error", func(t *testing.T) {
		mServ.EXPECT().GetShortURLFromCtx(gomock.Any()).Return(testShortURL, nil)
		mServ.EXPECT().GetOrigURLByShort(gomock.Any(), testShortURL).Return(models.OrigURL(""), errTest)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
