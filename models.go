package main

import "time"

type Quote struct {
	ID        int
	Text      string
	Author    string
	Tag       string
	CreatedAt time.Time
}
