package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type shortenBatchServicer interface {
	BatchShortenURL(context.Context, []models.OrigURL) ([]models.URLPair, error)
}

type ShortenBatchHandler struct {
	shortenBatchService shortenBatchServicer
	baseAddr            string
}

func NewShortenBatchHandler(shortenBatchService shortenBatchServicer, baseAddr string) *ShortenBatchHandler {
	return &ShortenBatchHandler{
		shortenBatchService: shortenBatchService,
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

func (h *ShortenBatchHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	var reqBody []shortenBatchReq
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var origs = make([]models.OrigURL, len(reqBody))
	for i, sbreq := range reqBody {
		origs[i] = models.OrigURL(sbreq.OrigURL)
	}

	pairs, err := h.shortenBatchService.BatchShortenURL(req.Context(), origs)
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
	err = json.NewEncoder(res).Encode(resBody)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
}
