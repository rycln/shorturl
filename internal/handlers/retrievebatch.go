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

type retrieveBatchServicer interface {
	GetUserURLs(context.Context, models.UserID) ([]models.URLPair, error)
}

type retrieveBatchAuthServicer interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// RetrieveBatchHandler handles requests to retrieve all URLs shortened by a user.
//
// Provides authenticated access to user's URL history in JSON format.
// The handler:
// 1. Extracts user ID from request context (set by auth middleware)
// 2. Fetches all user's URL pairs from storage
// 3. Returns formatted JSON response or 204 if no URLs exist
//
// Response codes:
//   - 200 OK: URLs found and returned
//   - 204 No Content: no URLs found for user
//   - 500 Internal Server Error: processing failure
type RetrieveBatchHandler struct {
	retrieveBatchService retrieveBatchServicer
	authService          retrieveBatchAuthServicer
	baseAddr             string
}

type errRetrieveBatchNotExist interface {
	error
	IsErrNotExist() bool
}

// NewRetrieveBatchHandler creates new user URLs handler instance.
func NewRetrieveBatchHandler(retrieveBatchService retrieveBatchServicer, authService retrieveBatchAuthServicer, baseAddr string) *RetrieveBatchHandler {
	return &RetrieveBatchHandler{
		retrieveBatchService: retrieveBatchService,
		authService:          authService,
		baseAddr:             baseAddr,
	}
}

type retBatchRes struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

// ServeHTTP implements http.Handler interface for user URLs endpoint.
//
// Expected request format:
//
//	GET /api/user/urls
//	Authorization: Bearer <token>
func (h *RetrieveBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uid, err := h.authService.GetUserIDFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	pairs, err := h.retrieveBatchService.GetUserURLs(req.Context(), models.UserID(uid))
	if e, ok := err.(errRetrieveBatchNotExist); ok && e.IsErrNotExist() {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var resBatch = make([]retBatchRes, len(pairs))
	for i, pair := range pairs {
		resBatch[i] = retBatchRes{
			ShortURL: h.baseAddr + "/" + string(pair.Short),
			OrigURL:  string(pair.Orig),
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&resBatch)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
}
