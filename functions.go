package main

import (
	"time"
)

func parseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

func getAllQuotes() []Quote {
	return quotes
}

func getQuotesByTag(tag string) []Quote {
	var newlist []Quote
	for i := range quotes {

		if quotes[i].Tag == tag {
			newlist = append(newlist, quotes[i])
		}

	}
	return newlist
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
