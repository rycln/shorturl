// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"context"
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

func TestStatsHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStats := mocks.NewMockstatsServicer(ctrl)

	statsHandler := NewStatsHandler(mStats)

	testStats := &models.Stats{
		URLs:  1,
		Users: 1,
	}

	testStatsJSON, err := json.Marshal(testStats)
	require.NoError(t, err)

	t.Run("valid test", func(t *testing.T) {
		mStats.EXPECT().GetStatsFromStorage(context.Background()).Return(testStats, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		statsHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(testStatsJSON)+"\n", string(resBody))
	})

	t.Run("service error", func(t *testing.T) {
		mStats.EXPECT().GetStatsFromStorage(context.Background()).Return(nil, errTest)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		statsHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
