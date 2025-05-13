package handlers

import (
	"bytes"
	"encoding/json"
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

func TestDeleteBatchHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)

	mProc := mocks.NewMockdeletionProcessor(ctrl)
	mAuth := mocks.NewMockdeleteBatchAuthServicer(ctrl)

	delBatchHandler := NewDeleteBatchHandler(mProc, mAuth)

	var testShortURLs = []models.ShortURL{
		"1",
		"2",
		"3",
	}
	t.Run("valid test", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testUserID, nil)
		mProc.EXPECT().AddURLsIntoDeletionQueue(gomock.Any(), gomock.Any())

		jsonReq, err := json.Marshal(&testShortURLs)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		delBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusAccepted, res.StatusCode)
	})

	t.Run("user id error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(models.UserID(""), errTest)

		jsonReq, err := json.Marshal(&testShortURLs)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		delBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("wrong json", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testUserID, nil)

		reqBody := strings.NewReader("not json")
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		delBatchHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
