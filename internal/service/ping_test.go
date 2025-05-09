package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPing_PingStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockPingStorage(ctrl)

	s := NewPing(mStrg)

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().Ping(context.Background()).Return(nil)

		err := s.PingStorage(context.Background())
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().Ping(context.Background()).Return(errTest)

		err := s.PingStorage(context.Background())
		assert.Error(t, err)
	})
}
