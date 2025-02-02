package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rycln/shorturl/internal/app/mem"
	"github.com/rycln/shorturl/internal/app/myhash"
)

type HandlerVariables struct {
	store mem.Storager
}

func NewHandlerVariables(store mem.Storager) HandlerVariables {
	return HandlerVariables{
		store: store,
	}
}

func (hv HandlerVariables) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fullURL := string(body)
	shortURL := myhash.Base62(fullURL)
	hv.store.AddURL(shortURL, fullURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", shortURL)))
}

func (hv HandlerVariables) ReturnURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimLeft(r.URL.Path, "/")
	fullURL, err := hv.store.GetURL(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
