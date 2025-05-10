package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type retrieveServicer interface {
	GetOrigURLByShort(context.Context, models.ShortURL) (models.OrigURL, error)
}

type RetrieveHandler struct {
	retrieveService retrieveServicer
}

type errRetrieveDeletedURL interface {
	error
	IsErrDeletedURL() bool
}

func NewRetrieveHandler(retrieveService retrieveServicer) *RetrieveHandler {
	return &RetrieveHandler{
		retrieveService: retrieveService,
	}
}

func (h *RetrieveHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	shortURL := req.Context().Value("short").(string)
	if shortURL == "" {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(errors.New("short URL empty value")))
		return
	}

	origURL, err := h.retrieveService.GetOrigURLByShort(req.Context(), models.ShortURL(shortURL))
	if e, ok := err.(errRetrieveDeletedURL); ok && e.IsErrDeletedURL() {
		res.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	res.Header().Set("Location", string(origURL))
	res.WriteHeader(http.StatusTemporaryRedirect)
}
