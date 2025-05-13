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

type PingHandler struct {
	pingService pingServicer
}

func NewPingHandler(pingService pingServicer) *PingHandler {
	return &PingHandler{
		pingService: pingService,
	}
}

func (h *PingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := h.pingService.PingStorage(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
	res.WriteHeader(http.StatusOK)
}
