package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "quotes.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/quotes/", quotesListDispatcher)
	mux.HandleFunc("POST /api/quotes/", createQuoteHandler)

	mux.HandleFunc("GET /api/quotes/{id}", quotesByIdHandler)
	mux.HandleFunc("PUT /api/quotes/{id}", updateQuoteByIdHandler)
	mux.HandleFunc("DELETE /api/quotes/{id}", deleteQuoteByIdHandler)

	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}

}
