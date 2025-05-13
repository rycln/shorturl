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

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type shortenServicer interface {
	ShortenURL(context.Context, models.UserID, models.OrigURL) (*models.URLPair, error)
}

type shortenAuthServicer interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

type ShortenHandler struct {
	shortenService shortenServicer
	authService    shortenAuthServicer
	baseAddr       string
}

type errShortenConflict interface {
	error
	IsErrConflict() bool
}

func NewShortenHandler(shortenService shortenServicer, authService shortenAuthServicer, baseAddr string) *ShortenHandler {
	return &ShortenHandler{
		shortenService: shortenService,
		authService:    authService,
		baseAddr:       baseAddr,
	}
}

func (h *ShortenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uid, err := h.authService.GetUserIDFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

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

	pair, err := h.shortenService.ShortenURL(req.Context(), uid, models.OrigURL(body))
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
