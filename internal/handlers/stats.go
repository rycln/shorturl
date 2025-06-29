// Package handlers provides HTTP request handlers for the application.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// StatsServicer defines the interface for statistics collection service.
type statsServicer interface {
	// GetStatsFromStorage retrieves the current service statistics from storage.
	GetStatsFromStorage(ctx context.Context) (*models.Stats, error)
}

// StatsHandler handles HTTP requests for service statistics.
type StatsHandler struct {
	statsService statsServicer
}

// NewStatsHandler creates a new StatsHandler instance.
func NewStatsHandler(statsService statsServicer) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetStats handles GET requests for service statistics.
func (h *StatsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	stats, err := h.statsService.GetStatsFromStorage(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(stats)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
	}
}
