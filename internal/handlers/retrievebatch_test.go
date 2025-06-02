package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/handlers/mocks"
	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveBatchHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mServ := mocks.NewMockretrieveBatchServicer(ctrl)
	mAuth := mocks.NewMockretrieveBatchAuthServicer(ctrl)

	retrieveBatchHandler := NewRetrieveBatchHandler(mServ, mAuth, testBaseAddr)

	testPairBatch := []models.URLPair{
		testPair,
	}

	resBatch := []retBatchRes{
		{
			OrigURL:  string(testPair.Orig),
			ShortURL: testBaseAddr + "/" + string(testPair.Short),
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mServ.EXPECT().GetUserURLs(gomock.Any(), testPair.UID).Return(testPairBatch, nil)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		jsonRes, err := json.Marshal(&resBatch)
		require.NoError(t, err)

		assert.Equal(t, string(jsonRes)+"\n", string(resBody))
	})

	t.Run("user id error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(models.UserID(""), errTest)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("no content error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mErr := mocks.NewMockerrRetrieveBatchNotExist(ctrl)
		mErr.EXPECT().IsErrNotExist().Return(true)
		mServ.EXPECT().GetUserURLs(gomock.Any(), testPair.UID).Return(nil, mErr)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("some service error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mServ.EXPECT().GetUserURLs(gomock.Any(), testPair.UID).Return(nil, errTest)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		retrieveBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
