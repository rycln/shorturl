package handlers

import (
	"context"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type pingServicer interface {
	PingStorage(context.Context) error
}

// PingHandler implements health check endpoint for database connectivity.
//
// Provides simple way to verify storage availability before processing requests.
// Handler flow:
// 1. Attempts to ping configured database
// 2. Returns status code based on connection check:
//   - 200 OK: storage is available
//   - 500 Internal Server Error: connection failed
type PingHandler struct {
	pingService pingServicer
}

// NewPingHandler creates new ping handler instance.
func NewPingHandler(pingService pingServicer) *PingHandler {
	return &PingHandler{
		pingService: pingService,
	}
}

// ServeHTTP implements http.Handler interface for ping endpoint.
//
// Expected request format:
//
//	GET /ping
func (h *PingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := h.pingService.PingStorage(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
	res.WriteHeader(http.StatusOK)
}
