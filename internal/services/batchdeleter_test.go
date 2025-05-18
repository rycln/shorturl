package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBatchDelete_DeleteURLsBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockBatchDeleterStorage(ctrl)

	s := NewBatchDeleter(mStrg)

	testDelReqs := []*models.DelURLReq{
		&testDelReq,
	}

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().DeleteRequestedURLs(context.Background(), testDelReqs).Return(nil)

		err := s.DeleteURLsBatch(context.Background(), testDelReqs)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().DeleteRequestedURLs(context.Background(), testDelReqs).Return(errTest)

		err := s.DeleteURLsBatch(context.Background(), testDelReqs)
		assert.Error(t, err)
	})
}
