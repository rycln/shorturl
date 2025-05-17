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

// RetrieveHandler handles requests to resolve shortened URLs.
//
// Implements HTTP redirection flow:
// 1. Extracts short URL ID from path parameter
// 2. Looks up original URL in storage
// 3. Returns 307 Redirect with original URL
//
// Response codes:
//   - 307 Temporary Redirect: successful lookup
//   - 410 Gone: URL was deleted
//   - 500 Internal Server Error: processing failure
type RetrieveHandler struct {
	retrieveService retrieveServicer
}

type errRetrieveDeletedURL interface {
	error
	IsErrDeletedURL() bool
}

// NewRetrieveHandler creates new redirect handler instance.
func NewRetrieveHandler(retrieveService retrieveServicer) *RetrieveHandler {
	return &RetrieveHandler{
		retrieveService: retrieveService,
	}
}

// ServeHTTP implements http.Handler interface for redirect endpoint.
//
// Expected request format:
//
//	GET /{id}
//	Content-Type: text/plain
func (h *RetrieveHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
