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

type DeleteBatchHandler struct {
	delProc     deletionProcessor
	authService deleteBatchAuthServicer
}

func NewDeleteBatchHandler(delProc deletionProcessor, authService deleteBatchAuthServicer) *DeleteBatchHandler {
	return &DeleteBatchHandler{
		delProc:     delProc,
		authService: authService,
	}
}

func (h *DeleteBatchHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
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
