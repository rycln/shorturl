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

type shortenBatchServicer interface {
	BatchShortenURL(context.Context, models.UserID, []models.OrigURL) ([]models.URLPair, error)
}

type shortenBatchAuthServicer interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// ShortenBatchHandler handles batch URL shortening requests.
//
// Processes multiple URLs in single operation while preserving order.
// Response maintains the same correlation IDs as in request for client-side matching.
//
// Response codes:
//   - 201 Created: all URLs processed successfully
//   - 400 Bad Request: invalid input data
//   - 500 Internal Server Error: processing failure
type ShortenBatchHandler struct {
	shortenBatchService shortenBatchServicer
	authService         shortenBatchAuthServicer
	baseAddr            string
}

// NewShortenBatchHandler creates new batch handler instance.
func NewShortenBatchHandler(shortenBatchService shortenBatchServicer, authService shortenBatchAuthServicer, baseAddr string) *ShortenBatchHandler {
	return &ShortenBatchHandler{
		shortenBatchService: shortenBatchService,
		authService:         authService,
		baseAddr:            baseAddr,
	}
}

type shortenBatchReq struct {
	ID      string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}

type shortenBatchRes struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// ServeHTTP implements http.Handler interface for batch endpoint.
//
// Expected request format:
//
//	POST /api/shorten/batch
//	Content-Type: application/json
//	Authorization: Bearer <token>
func (h *ShortenBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uid, err := h.authService.GetUserIDFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var reqBody []shortenBatchReq
	err = json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var origs = make([]models.OrigURL, len(reqBody))
	for i, sbreq := range reqBody {
		origs[i] = models.OrigURL(sbreq.OrigURL)
	}

	pairs, err := h.shortenBatchService.BatchShortenURL(req.Context(), uid, origs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var resBody = make([]shortenBatchRes, len(pairs))
	for i, pair := range pairs {
		resBody[i] = shortenBatchRes{
			ID:       reqBody[i].ID,
			ShortURL: h.baseAddr + "/" + string(pair.Short),
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(res).Encode(resBody)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
}
