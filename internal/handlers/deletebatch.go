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

type deletionProcessor interface {
	AddURLsIntoDeletionQueue(models.UserID, []models.ShortURL)
}

type deleteBatchAuthServicer interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// DeleteBatchHandler handles asynchronous batch URL deletion requests.
//
// The handler:
// 1. Extracts user ID from request context (set by auth middleware)
// 2. Queues deletion tasks in background
// 3. Immediately returns 202 Accepted
//
// Response codes:
//   - 202 Accepted: request queued for processing
//   - 500 Internal Server Error: queue failure
//
// Only URL owner can successfully delete URLs.
type DeleteBatchHandler struct {
	delProc     deletionProcessor
	authService deleteBatchAuthServicer
}

// NewDeleteBatchHandler creates new batch deletion handler instance.
func NewDeleteBatchHandler(delProc deletionProcessor, authService deleteBatchAuthServicer) *DeleteBatchHandler {
	return &DeleteBatchHandler{
		delProc:     delProc,
		authService: authService,
	}
}

// ServeHTTP implements http.Handler interface for batch deletion endpoint.
//
// Expected request format:
//
//	DELETE /api/user/urls
//	Content-Type: application/json
//	Authorization: Bearer <token>
func (h *DeleteBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uid, err := h.authService.GetUserIDFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var surls []models.ShortURL
	err = json.NewDecoder(req.Body).Decode(&surls)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	h.delProc.AddURLsIntoDeletionQueue(models.UserID(uid), surls)
	res.WriteHeader(http.StatusAccepted)
}
