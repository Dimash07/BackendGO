package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := getAllQuotes()
	json.NewEncoder(w).Encode(response)
}

func quotesByTagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tag := r.URL.Query().Get("tag")
	fmt.Println(tag)
	response := getQuotesByTag(tag)

	json.NewEncoder(w).Encode(response)
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

func quotesListDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("tag") != "" {
		quotesByTagHandler(w, r)
		return
	}
	quotesHandler(w, r)
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
