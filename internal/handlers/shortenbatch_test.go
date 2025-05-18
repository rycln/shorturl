package handlers

import (
	"bytes"
	"encoding/json"
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

func TestShortenBatchHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mShort := mocks.NewMockshortenBatchServicer(ctrl)
	mAuth := mocks.NewMockshortenBatchAuthServicer(ctrl)

	shortenBatchHandler := NewShortenBatchHandler(mShort, mAuth, testBaseAddr)

	testPairBatch := []models.URLPair{
		testPair,
	}

	reqBatch := []shortenBatchReq{
		{
			ID:      string(testPair.UID),
			OrigURL: string(testPair.Orig),
		},
	}

	resBatch := []shortenBatchRes{
		{
			ID:       string(testPair.UID),
			ShortURL: testBaseAddr + "/" + string(testPair.Short),
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mShort.EXPECT().BatchShortenURL(gomock.Any(), testPair.UID, gomock.Any()).Return(testPairBatch, nil)

		jsonReq, err := json.Marshal(&reqBatch)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		jsonRes, err := json.Marshal(&resBatch)
		require.NoError(t, err)

		assert.Equal(t, string(jsonRes)+"\n", string(resBody))
	})

	t.Run("user id error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(models.UserID(""), errTest)

		jsonReq, err := json.Marshal(&reqBatch)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("wrong request", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)

		reqBody := strings.NewReader("wrong json")
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("some service error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mShort.EXPECT().BatchShortenURL(gomock.Any(), testPair.UID, gomock.Any()).Return(nil, errTest)

		jsonReq, err := json.Marshal(&reqBatch)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		shortenBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

}
