package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /api/quotes/", quotesListDispatcher)
	mux.HandleFunc("POST /api/quotes/", createQuoteHandler)

	mux.HandleFunc("GET /api/quotes/{id}", quotesByIdHandler)
	mux.HandleFunc("PUT /api/quotes/{id}", updateQuoteByIdHandler)
	mux.HandleFunc("DELETE /api/quotes/{id}", deleteQuoteByIdHandler)

	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}
