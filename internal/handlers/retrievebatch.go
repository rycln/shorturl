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

type RetrieveBatchHandler struct {
	retrieveBatchService retrieveBatchServicer
	authService          retrieveBatchAuthServicer
	baseAddr             string
}

type errRetrieveBatchNotExist interface {
	error
	IsErrNotExist() bool
}

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
