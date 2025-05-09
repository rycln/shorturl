package storage

import (
	"errors"

	"github.com/rycln/shorturl/internal/models"
)

const (
	testUserID       models.UserID   = "1"
	testShortURL     models.ShortURL = "abc123"
	testDeletedShort models.ShortURL = "321cba"
	testOrigURL      models.OrigURL  = "https://practicum.yandex.ru/"
	testFileName                     = "test"
)

var (
	errTest = errors.New("test error")

	testPair = models.URLPair{
		UID:   testUserID,
		Short: testShortURL,
		Orig:  testOrigURL,
	}

	testDelReq = models.DelURLReq{
		UID:   testUserID,
		Short: testShortURL,
	}
)
