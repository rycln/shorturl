package handlers

import (
	"errors"
	"time"
)

const (
	testBaseAddr        = "http://localhost:8080"
	testHashVal         = "abc"
	testTimeoutDuration = time.Duration(2) * time.Second
	testKey             = "test_key"
	testID              = "1"
)

var errTest = errors.New("test error")

func testHash(str string) string {
	return testHashVal
}
