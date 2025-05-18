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

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	fmt.Println("Status:", rr.Code)

	fmt.Println("Response:", rr.Body.String())

	// Output:
	// Status: 201
	// Response: http://localhost:8080/abc
}
