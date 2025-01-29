package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rycln/shorturl/internal/app/mem"
)

var m = mem.MemStorage{
	Storage: make(map[string]string),
}

func shortenURL(s mem.Memorizer, url string) string {
	return s.AddURL(url)
}

func retrieveURL(s mem.Memorizer, shortURL string) (string, error) {
	return s.GetURL(shortURL)
}

func shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	sURL := shortenURL(&m, string(body))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", sURL)))
}

func urlReturn(w http.ResponseWriter, r *http.Request) {
	sURL := strings.TrimLeft(r.URL.Path, "/")
	fullURL, err := retrieveURL(&m, sURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", shortener)
	mux.HandleFunc("/{id}", urlReturn)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
