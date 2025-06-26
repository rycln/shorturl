package services

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// StatsStorage defines the interface for storage operations related to statistics.
// Implementations of this interface should provide methods to retrieve URL and user statistics.
type StatsStorage interface {
	// GetStats retrieves aggregated statistics (URLs and users count) from storage.
	GetStats(context.Context) (*models.Stats, error)
}

// StatsCollector provides methods for collecting and retrieving service statistics.
// It acts as an intermediary between the HTTP layer and the storage layer.
type StatsCollector struct {
	strg StatsStorage
}

// NewStatsCollector creates a new instance of StatsCollector with the given storage.
func NewStatsCollector(strg StatsStorage) *StatsCollector {
	return &StatsCollector{
		strg: strg,
	}
}

// GetStatsFromStorage retrieves the latest statistics from the underlying storage.
func (s *StatsCollector) GetStatsFromStorage(ctx context.Context) (*models.Stats, error) {
	stats, err := s.strg.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
