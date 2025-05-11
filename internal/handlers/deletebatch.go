package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type deleteBatchServicer interface {
	UserURLsAsyncDeletion(models.UserID, []models.ShortURL)
}

type DeleteBatchHandler struct {
	deleteBatchService deleteBatchServicer
}

func NewDeleteBatchHandler(deleteBatchService deleteBatchServicer) *DeleteBatchHandler {
	return &DeleteBatchHandler{
		deleteBatchService: deleteBatchService,
	}
}

func (h *DeleteBatchHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	uid := req.Context().Value("uid").(string)

	var surls []models.ShortURL
	err := json.NewDecoder(req.Body).Decode(&surls)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	h.deleteBatchService.UserURLsAsyncDeletion(models.UserID(uid), surls)
	res.WriteHeader(http.StatusAccepted)
}
