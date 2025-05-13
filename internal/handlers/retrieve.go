package handlers

import (
	"context"
	"net/http"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type retrieveServicer interface {
	GetShortURLFromCtx(context.Context) (models.ShortURL, error)
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
	shortURL, err := h.retrieveService.GetShortURLFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	origURL, err := h.retrieveService.GetOrigURLByShort(req.Context(), shortURL)
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
