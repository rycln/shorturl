package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStatsCollector_GetStatsFromStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockStatsStorage(ctrl)

	statcollector := NewStatsCollector(mStrg)

	testStats := &models.Stats{
		URLs:  1,
		Users: 1,
	}

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().GetStats(context.Background()).Return(testStats, nil)

		stats, err := statcollector.GetStatsFromStorage(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, testStats, stats)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetStats(context.Background()).Return(nil, errTest)

		_, err := statcollector.GetStatsFromStorage(context.Background())
		assert.Error(t, err)
	})
}
