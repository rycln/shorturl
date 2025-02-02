package main

import (
	"net/http"

	"github.com/rycln/shorturl/internal/app/handlers"
	"github.com/rycln/shorturl/internal/app/mem"
)

func main() {
	store := mem.NewSimpleMemStorage()
	hv := handlers.NewHandlerVariables(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hv.ShortenURL)
	mux.HandleFunc("/{id}", hv.ReturnURL)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
