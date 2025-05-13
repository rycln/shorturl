package handlers

import (
	"errors"

	"github.com/rycln/shorturl/internal/models"
)

const (
	testBaseAddr                     = "test/"
	testUserID       models.UserID   = "1"
	testShortURL     models.ShortURL = "abc123"
	testDeletedShort models.ShortURL = "321cba"
	testOrigURL      models.OrigURL  = "https://practicum.yandex.ru/"
)

var (
	errTest = errors.New("test error")

	testPair = models.URLPair{
		UID:   testUserID,
		Short: testShortURL,
		Orig:  testOrigURL,
	}
)
