package models

import "time"

type ChapterQueryParams struct {
	Limit              int      `json:"limit"`
	Offset             int      `json:"offset"`
	Ids                []string `json:"ids"`
	MangaId            string   `json:"manga"`
	TranslatedLanguage []string `json:"translatedLanguage"`
}

type Chapter struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Attributes    ChapterAttributes     `json:"attributes"`
	Relationships []ChapterRelationship `json:"relationships"`
}

type ChapterAttributes struct {
	Volume             string    `json:"volume"`
	Chapter            string    `json:"chapter"`
	Title              string    `json:"title"`
	TranslatedLanguage string    `json:"translatedLanguage"`
	ExternalURL        *string   `json:"externalUrl"`
	IsUnavailable      bool      `json:"isUnavailable"`
	PublishAt          time.Time `json:"publishAt"`
	ReadableAt         time.Time `json:"readableAt"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Pages              int       `json:"pages"`
	Version            int       `json:"version"`
}

type ChapterRelationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type ChapterListResponse struct {
	Data     []Chapter `json:"data"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Total    int       `json:"total"`
	Result   string    `json:"result"`
	Response string    `json:"response"`
}
