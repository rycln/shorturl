package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

type apiShortenServicer interface {
	ShortenURL(context.Context, models.OrigURL) (*models.URLPair, error)
}

type APIShortenHandler struct {
	apiShortenService shortenServicer
	baseAddr          string
}

type errAPIShortenConflict interface {
	error
	IsErrConflict() bool
}

func NewAPIShortenHandler(apiShortenService apiShortenServicer, baseAddr string) *APIShortenHandler {
	return &APIShortenHandler{
		apiShortenService: apiShortenService,
		baseAddr:          baseAddr,
	}
}

type apiShortenReq struct {
	URL string `json:"url"`
}

type apiShortenRes struct {
	Result string `json:"result"`
}

func (h *APIShortenHandler) HandleHTTP(res http.ResponseWriter, req *http.Request) {
	var reqBody apiShortenReq
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
	_, err = url.ParseRequestURI(reqBody.URL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	pair, err := h.apiShortenService.ShortenURL(req.Context(), models.OrigURL(reqBody.URL))
	if e, ok := err.(errAPIShortenConflict); ok && e.IsErrConflict() {
		err := h.sendResponse(res, http.StatusConflict, string(pair.Short))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
			return
		}
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	err = h.sendResponse(res, http.StatusCreated, string(pair.Short))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}
}

func (h *APIShortenHandler) sendResponse(res http.ResponseWriter, code int, shortURL string) error {
	var resBody apiShortenRes
	resBody.Result = h.baseAddr + shortURL

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	err := json.NewEncoder(res).Encode(resBody)
	if err != nil {
		return err
	}
	return nil
}
