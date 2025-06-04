package models

import "time"

// Chapter represents a manga chapter from MangaDex
type Chapter struct {
	ID        string
	Title     string
	Volume    string
	Chapter   string
	Pages     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
