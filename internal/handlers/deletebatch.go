package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type deletionProcessor interface {
	AddURLsIntoDeletionQueue(models.UserID, []models.ShortURL)
}

type DeleteBatchHandler struct {
	delProc deletionProcessor
}

func NewDeleteBatchHandler(delProc deletionProcessor) *DeleteBatchHandler {
	return &DeleteBatchHandler{
		delProc: delProc,
	}
}

func (h *DeleteBatchHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	uid, ok := req.Context().Value(contextkeys.UserID).(models.UserID)
	if !ok {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(errors.New("short URL value is empty")))
		return
	}

	var surls []models.ShortURL
	err := json.NewDecoder(req.Body).Decode(&surls)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	h.delProc.AddURLsIntoDeletionQueue(models.UserID(uid), surls)
	res.WriteHeader(http.StatusAccepted)
}
