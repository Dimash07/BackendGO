package main

import (
	"database/sql"
	"io"
	"log/slog"
	"os"
	"time"
)

func parseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

func getAllQuotes(db *sql.DB) []Quote {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(multi, nil))
	if db == nil {
		logger.Info("DB is nil!")
	}
	rows, err := db.Query("SELECT id, text, author, tag, created_at FROM quotes")
	if err != nil {
		logger.Info("Ошибка при запросе к БД:", err)
		return []Quote{}
	}
	defer rows.Close()

	var list []Quote
	for rows.Next() {
		var q Quote
		var CreatedAtStr string
		err := rows.Scan(&q.ID, &q.Text, &q.Author, &q.Tag, &CreatedAtStr)
		if err != nil {
			logger.Info("Ошибка скан строка:", err)
			continue
		}
		q.CreatedAt, _ = time.Parse(time.RFC3339, CreatedAtStr)
		list = append(list, q)
	}

	return list
}

func getQuotesByTag(db *sql.DB, tag string) ([]Quote, error) {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(multi, nil))
	rows, err := db.Query("SELECT id, text, author, tag, created_at FROM quotes WHERE tag = ?", tag)
	if err != nil {
		logger.Error("Search by tag failed", "tag", tag, "error", err)
		return nil, err
	}
	defer rows.Close()

	var list []Quote
	for rows.Next() {
		var q Quote
		var createdAtStr string
		if err := rows.Scan(&q.ID, &q.Text, &q.Author, &q.Tag, &createdAtStr); err != nil {
			continue
		}
		q.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, q)
	}
	return list, nil
}

func getQuoteById(id int) *Quote {

	for i := range quotes {
		if quotes[i].ID == id {
			return &quotes[i]

		}
	}
	return nil
}

func updateQuoteById(id int, text string, author string, tag string, createdat time.Time) *Quote {

	for i := range quotes {
		if quotes[i].ID == id {
			quotes[i].Text = text
			quotes[i].Author = author
			quotes[i].Tag = tag
			quotes[i].CreatedAt = createdat

			return &quotes[i]
		}

	}
	return nil

}

func deleteQuoteById(id int) bool {

	for i, q := range quotes {

		if q.ID == id {
			quotes = append(quotes[:i], quotes[i+1:]...)
			return true
		}
	}
	return false
}

var LastId = 0

func createQuote(text, author, tag string, timestamp time.Time) *Quote {

	maxID := 0
	for _, q := range quotes {
		if q.ID > maxID {
			maxID = q.ID
		}
	}
	newQuote := Quote{
		ID:        maxID + 1,
		Text:      text,
		Author:    author,
		Tag:       tag,
		CreatedAt: timestamp,
	}
	quotes = append(quotes, newQuote)
	return &quotes[len(quotes)-1]
}
