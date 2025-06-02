package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPing := mocks.NewMockpingServicer(ctrl)

	pingHandler := NewPingHandler(mPing)

	t.Run("valid test", func(t *testing.T) {
		mPing.EXPECT().PingStorage(gomock.Any()).Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		pingHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("some service error", func(t *testing.T) {
		mPing.EXPECT().PingStorage(gomock.Any()).Return(errTest)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		pingHandler.ServeHTTP(w, req)

		res := w.Result()
		defer func() {
			err := res.Body.Close()
			require.NoError(t, err)
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
