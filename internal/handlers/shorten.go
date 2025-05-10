package handlers

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type shortenServicer interface {
	ShortenURL(context.Context, models.OrigURL) (*models.URLPair, error)
}

type ShortenHandler struct {
	shortenService shortenServicer
	baseAddr       string
}

type errShortenConflict interface {
	error
	IsErrConflict() bool
}

func NewShortenHandler(shortenService shortenServicer, baseAddr string) *ShortenHandler {
	return &ShortenHandler{
		shortenService: shortenService,
		baseAddr:       baseAddr,
	}
}

func (h *ShortenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	pair, err := h.shortenService.ShortenURL(req.Context(), models.OrigURL(body))
	if e, ok := err.(errShortenConflict); ok && e.IsErrConflict() {
		h.sendResponse(res, http.StatusConflict, string(pair.Short))
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	h.sendResponse(res, http.StatusCreated, string(pair.Short))
}

func (h *ShortenHandler) sendResponse(res http.ResponseWriter, code int, shortURL string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(code)
	res.Write([]byte(h.baseAddr + "/" + shortURL))
}
