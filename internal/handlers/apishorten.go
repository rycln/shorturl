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

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type apiShortenServicer interface {
	ShortenURL(context.Context, models.UserID, models.OrigURL) (*models.URLPair, error)
}

type apiShortenAuthServicer interface {
	GetUserIDFromCtx(context.Context) (models.UserID, error)
}

// APIShortenHandler handles URL shortening requests.
//
// The handler:
// 1. Extracts user ID from request context (set by auth middleware)
// 2. Validates input URL from request body
// 3. Processes through shortening service
// 4. Returns appropriate HTTP response and body:
//   - 201 Created: successful shortening
//   - 400 Bad Request: invalid input
//   - 409 Conflict: URL already exists
//   - 500 Internal Server Error: processing failure
type APIShortenHandler struct {
	apiShortenService apiShortenServicer
	authService       apiShortenAuthServicer
	baseAddr          string
}

type errAPIShortenConflict interface {
	error
	IsErrConflict() bool
}

// NewAPIShortenHandler creates a new handler instance with required dependencies.
func NewAPIShortenHandler(apiShortenService apiShortenServicer, authService apiShortenAuthServicer, baseAddr string) *APIShortenHandler {
	return &APIShortenHandler{
		apiShortenService: apiShortenService,
		authService:       authService,
		baseAddr:          baseAddr,
	}
}

type apiShortenReq struct {
	URL string `json:"url"`
}

type apiShortenRes struct {
	Result string `json:"result"`
}

// ServeHTTP implements http.Handler interface for the endpoint.
//
// Expected request format:
//
//	POST /api/shorten
//	Content-Type: application/json
//	Authorization: Bearer <token>
func (h *APIShortenHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	uid, err := h.authService.GetUserIDFromCtx(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Log.Debug("path:"+req.URL.Path, zap.Error(err))
		return
	}

	var reqBody apiShortenReq
	err = json.NewDecoder(req.Body).Decode(&reqBody)
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

	pair, err := h.apiShortenService.ShortenURL(req.Context(), uid, models.OrigURL(reqBody.URL))
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
	resBody.Result = h.baseAddr + "/" + shortURL

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	err := json.NewEncoder(res).Encode(resBody)
	if err != nil {
		return err
	}
	return nil
}
