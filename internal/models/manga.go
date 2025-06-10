package models

type Manga struct {
	ID            string
	Title         string
	Description   string
	Status        string
	Year          int
	Tags          []Tag
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
}

type Tag struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        map[string]string `json:"name"`
		Description map[string]string `json:"description"`
		Group       string            `json:"group"`
		Version     int               `json:"version"`
	} `json:"attributes"`
}

type MangaListResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
	Data     []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title         map[string]string `json:"title"`
			Description   map[string]string `json:"description"`
			Status        string            `json:"status"`
			Year          int               `json:"year"`
			IsLocked      bool              `json:"isLocked"`
			Links         map[string]string `json:"links"`
			ContentRating string            `json:"contentRating"`
			Tags          []Tag             `json:"tags"`
		} `json:"attributes"`
		Relationships []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"relationships"`
	} `json:"data"`
}
