package models

// Manga represents a manga from MangaDex
type Manga struct {
	ID          string
	Title       string
	Description string
	Status      string
	Year        int
	Tags        []Tag
	CoverURL    string
}

// Tag represents a manga tag
type Tag struct {
	ID   string
	Name string
}
