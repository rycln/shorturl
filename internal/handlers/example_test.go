package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/handlers/mocks"
	"github.com/rycln/shorturl/internal/models"
)

func ExampleShortenHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mShort := mocks.NewMockshortenServicer(ctrl)
	mAuth := mocks.NewMockshortenAuthServicer(ctrl)
	baseAddr := "http://localhost:8080"

	handler := NewShortenHandler(mShort, mAuth, baseAddr)

	pair := &models.URLPair{
		UID:   "user_1",
		Short: "abc",
		Orig:  "https://example.com",
	}
	mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(pair.UID, nil)
	mShort.EXPECT().ShortenURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(pair, nil)

	body := strings.NewReader(string(pair.Orig))
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "some.valid.jwt")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	fmt.Println("Response:", w.Body.String())

	// Output:
	// Status: 201
	// Response: http://localhost:8080/abc
}

func ExampleRetrieveHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mServ := mocks.NewMockretrieveServicer(ctrl)

	handler := NewRetrieveHandler(mServ)

	shortURL := models.ShortURL("abc123")
	origURL := models.OrigURL("https://example.com")
	mServ.EXPECT().GetShortURLFromCtx(gomock.Any()).Return(shortURL, nil)
	mServ.EXPECT().GetOrigURLByShort(gomock.Any(), shortURL).Return(origURL, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	fmt.Println("Location:", w.Header().Get("Location"))

	// Output:
	// Status: 307
	// Location: https://example.com
}

func ExampleAPIShortenHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mShort := mocks.NewMockapiShortenServicer(ctrl)
	mAuth := mocks.NewMockapiShortenAuthServicer(ctrl)
	baseAddr := "http://localhost:8080"

	handler := NewAPIShortenHandler(mShort, mAuth, baseAddr)

	pair := &models.URLPair{
		UID:   "user_1",
		Short: "abc",
		Orig:  "https://example.com",
	}
	mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(pair.UID, nil)
	mShort.EXPECT().ShortenURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(pair, nil)

	body := strings.NewReader(`{"url":"https://example.com"}`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "some.valid.jwt")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	fmt.Println("Response:", w.Body.String())

	// Output:
	// Status: 201
	// Response: {"result":"http://localhost:8080/abc"}
}

func ExamplePingHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mPing := mocks.NewMockpingServicer(ctrl)

	handler := NewPingHandler(mPing)

	mPing.EXPECT().PingStorage(gomock.Any()).Return(nil)

	req := httptest.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	// Output:
	// Status: 200
}

func ExampleShortenBatchHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mShort := mocks.NewMockshortenBatchServicer(ctrl)
	mAuth := mocks.NewMockshortenBatchAuthServicer(ctrl)
	baseAddr := "http://localhost:8080"

	handler := NewShortenBatchHandler(mShort, mAuth, baseAddr)

	pair := models.URLPair{
		UID:   "user_1",
		Short: "abc",
		Orig:  "https://example.com",
	}
	pairBatch := []models.URLPair{
		pair,
	}
	mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(pair.UID, nil)
	mShort.EXPECT().BatchShortenURL(gomock.Any(), pair.UID, gomock.Any()).Return(pairBatch, nil)

	body := strings.NewReader(`[{"correlation_id":"123","original_url":"https://example.com"}]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "some.valid.jwt")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	fmt.Println("Response:", w.Body.String())

	// Output:
	// Status: 201
	// Response: [{"correlation_id":"123","short_url":"http://localhost:8080/abc"}]
}

func ExampleRetrieveBatchHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mServ := mocks.NewMockretrieveBatchServicer(ctrl)
	mAuth := mocks.NewMockretrieveBatchAuthServicer(ctrl)
	baseAddr := "http://localhost:8080"

	handler := NewRetrieveBatchHandler(mServ, mAuth, baseAddr)

	pair := models.URLPair{
		UID:   "user_1",
		Short: "abc",
		Orig:  "https://example.com",
	}
	pairBatch := []models.URLPair{
		pair,
	}
	mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(pair.UID, nil)
	mServ.EXPECT().GetUserURLs(gomock.Any(), pair.UID).Return(pairBatch, nil)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "some.valid.jwt")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	fmt.Println("Response:", w.Body.String())

	// Output:
	// Status: 200
	// Response: [{"short_url":"http://localhost:8080/abc","original_url":"https://example.com"}]
}

func ExampleDeleteBatchHandler_ServeHTTP() {
	ctrl := gomock.NewController(nil)

	mProc := mocks.NewMockdeletionProcessor(ctrl)
	mAuth := mocks.NewMockdeleteBatchAuthServicer(ctrl)

	handler := NewDeleteBatchHandler(mProc, mAuth)

	userID := models.UserID("user1")
	mAuth.EXPECT().GetUserIDFromCtx(gomock.Any()).Return(userID, nil)
	mProc.EXPECT().AddURLsIntoDeletionQueue(gomock.Any(), gomock.Any())

	body := strings.NewReader(`["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Authorization", "some.valid.jwt")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	fmt.Println("Status:", w.Code)

	// Output:
	// Status: 202
}
