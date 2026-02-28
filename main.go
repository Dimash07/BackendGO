package main

import (
	"database/sql"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(multi, nil))

	logger.Info("Database init started")
	err := godotenv.Load()

	dbPath := os.Getenv("DATABASE_URL")
	logger.Info(dbPath)
	if dbPath == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	logger.Info("Database connected")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/quotes/", quotesListDispatcher(db))
	mux.HandleFunc("POST /api/quotes/", createQuoteHandler)

	mux.HandleFunc("GET /api/quotes/{id}", quotesByIdHandler)
	mux.HandleFunc("PUT /api/quotes/{id}", updateQuoteByIdHandler)
	mux.HandleFunc("DELETE /api/quotes/{id}", deleteQuoteByIdHandler)
	log.Println("Server started on :9090")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal("port failed", err)
	}

}
