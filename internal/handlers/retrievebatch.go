package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type retrieveBatchServicer interface {
	GetUserURLs(context.Context, models.UserID) ([]models.URLPair, error)
}

type RetrieveBatchHandler struct {
	retrieveBatchService retrieveBatchServicer
	baseAddr             string
}

type errRetrieveBatchNotExist interface {
	error
	IsErrNotExist() bool
}

func NewRetrieveBatchHandler(retrieveBatchService retrieveBatchServicer, baseAddr string) *RetrieveBatchHandler {
	return &RetrieveBatchHandler{
		retrieveBatchService: retrieveBatchService,
		baseAddr:             baseAddr,
	}
}

type retBatchRes struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

func (h *RetrieveBatchHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	uid, ok := req.Context().Value(contextkeys.UserID).(models.UserID)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(errors.New("short URL value is empty")))
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
