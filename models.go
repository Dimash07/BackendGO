package main

import "time"

type Quote struct {
    ID        int       `json:"id"`
    Text      string    `json:"text"`
    Author    string    `json:"author"`
    Tag       string    `json:"tag"`
    CreatedAt time.Time `json:"created_at"`
}
