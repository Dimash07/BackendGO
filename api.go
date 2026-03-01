package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func quotesHandler(db *sql.DB) http.HandlerFunc {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(multi, nil))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tag := r.URL.Query().Get("tag")

		var response []Quote
		var err error

		if tag != "" {
			// Если тег передан в URL (?tag=work)
			response, err = getQuotesByTag(db, tag)
		} else {
			// Если тега нет, вызываем твою старую функцию
			response = getAllQuotes(db)
		}

		if err != nil {
			logger.Error("Error fetching quotes", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(response)
	}
}

func quotesByTagHandler(db *sql.DB) http.HandlerFunc {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(multi, nil))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tag := r.URL.Query().Get("tag")
		if tag == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Tag parameter is required"})
			return
		}

		response, err := getQuotesByTag(db, tag)
		if err != nil {
			logger.Error("Failed to fetch quotes", "tag", tag, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if response == nil {
			response = []Quote{}
		}

		json.NewEncoder(w).Encode(response)
	}
}

func quotesByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idst := r.PathValue("id")

	id, err := strconv.Atoi(idst)

	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	response := getQuoteById(id)

	if response == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(response)

}

func updateQuoteByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idst := r.PathValue("id")

	id, err := strconv.Atoi(idst)

	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	var updatequote Quote
	err = json.NewDecoder(r.Body).Decode(&updatequote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if updatequote.Text == "" || updatequote.Author == "" {
		http.Error(w, "Text and Author are required fields", http.StatusBadRequest)
		return
	}

	response := updateQuoteById(
		id,
		updatequote.Text,
		updatequote.Author,
		updatequote.Tag,
		time.Now(),
	)
	if response == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func deleteQuoteByIdHandler(w http.ResponseWriter, r *http.Request) {
	idst := r.PathValue("id")
	id, err := strconv.Atoi(idst)

	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	isDeleted := deleteQuoteById(id)
	if isDeleted {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Quote not found", http.StatusNotFound)
	}
}

func createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newQuote Quote
	err := json.NewDecoder(r.Body).Decode(&newQuote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newQuote.Text == "" || newQuote.Author == "" {
		http.Error(w, "Text and Author are required fields", http.StatusBadRequest)
		return
	}
	response := createQuote(newQuote.Text, newQuote.Author, newQuote.Tag, time.Now())

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
