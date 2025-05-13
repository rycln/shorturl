package handlers

import (
	"bytes"
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

func TestAPIShortenHandler_HandleHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mShort := mocks.NewMockapiShortenServicer(ctrl)
	mAuth := mocks.NewMockapiShortenAuthServicer(ctrl)

	apiShortenHandler := NewAPIShortenHandler(mShort, mAuth, testBaseAddr)

	t.Run("valid test", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(testPair.UID, nil)
		mShort.EXPECT().ShortenURL(gomock.Any(), testPair.UID, testPair.Orig).Return(&testPair, nil)

		var reqOrig = apiShortenReq{
			URL: string(testOrigURL),
		}
		jsonReq, err := json.Marshal(&reqOrig)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonReq)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		apiShortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var resOrig = apiShortenRes{
			Result: testBaseAddr + "/" + string(testPair.Short),
		}
		jsonRes, err := json.Marshal(&resOrig)
		require.NoError(t, err)

		assert.Equal(t, string(jsonRes)+"\n", string(resBody))
	})

	t.Run("user id error", func(t *testing.T) {
		mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(models.UserID(""), errTest)

		var reqOrig = apiShortenReq{
			URL: string(testOrigURL),
		}
		jsonData, err := json.Marshal(&reqOrig)
		require.NoError(t, err)
		reqBody := bytes.NewReader(jsonData)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		w := httptest.NewRecorder()
		apiShortenHandler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
