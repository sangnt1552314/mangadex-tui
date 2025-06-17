package models

const (
	OrderByCreatedAt   = "createdAt"
	OrderByRating      = "rating"
	OrderByFollowCount = "followedCount"
)

type MangaQueryParams struct {
	Limit         int               `json:"limit"`
	ContentRating []string          `json:"contentRating"`
	Order         map[string]string `json:"order"`
	Includes      []string          `json:"includes"`
	HasChapters   bool              `json:"hasAvailableChapters"`
}

type Manga struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title                          map[string]string   `json:"title"`
		AltTitles                      []map[string]string `json:"altTitles"`
		Description                    map[string]string   `json:"description"`
		IsLocked                       bool                `json:"isLocked"`
		Links                          map[string]string   `json:"links"`
		OriginalLanguage               string              `json:"originalLanguage"`
		LastVolume                     string              `json:"lastVolume"`
		LastChapter                    string              `json:"lastChapter"`
		PublicationDemographic         string              `json:"publicationDemographic"`
		Status                         string              `json:"status"`
		Year                           int                 `json:"year"`
		ContentRating                  string              `json:"contentRating"`
		Tags                           []Tag               `json:"tags"`
		State                          string              `json:"state"`
		ChapterNumbersResetOnNewVolume bool                `json:"chapterNumbersResetOnNewVolume"`
		CreatedAt                      string              `json:"createdAt"`
		UpdatedAt                      string              `json:"updatedAt"`
		Version                        int                 `json:"version"`
		AvailableTranslatedLanguages   []string            `json:"availableTranslatedLanguages"`
		LatestUploadedChapter          string              `json:"latestUploadedChapter"`
	} `json:"attributes"`
	Relationships []MangaRelationship `json:"relationships"`
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
	Relationships []Relationship `json:"relationships"`
}

type MangaRelationship struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		// Common fields
		Name      string            `json:"name"`
		ImageUrl  *string           `json:"imageUrl"`
		Biography map[string]string `json:"biography"`
		CreatedAt string            `json:"createdAt"`
		UpdatedAt string            `json:"updatedAt"`
		Version   int               `json:"version"`

		// Social media links
		Twitter   *string `json:"twitter"`
		Pixiv     *string `json:"pixiv"`
		MelonBook *string `json:"melonBook"`
		FanBox    *string `json:"fanBox"`
		Booth     *string `json:"booth"`
		Namicomi  *string `json:"namicomi"`
		NicoVideo *string `json:"nicoVideo"`
		Skeb      *string `json:"skeb"`
		Fantia    *string `json:"fantia"`
		Tumblr    *string `json:"tumblr"`
		Youtube   *string `json:"youtube"`
		Weibo     *string `json:"weibo"`
		Naver     *string `json:"naver"`
		Website   *string `json:"website"`

		// Cover art specific fields
		Description string `json:"description"`
		Volume      string `json:"volume"`
		FileName    string `json:"fileName"`
		Locale      string `json:"locale"`
	} `json:"attributes"`
}

type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MangaListResponse struct {
	Result   string  `json:"result"`
	Response string  `json:"response"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
	Data     []Manga `json:"data"`
}

type CoverArt struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Description string `json:"description"`
		Volume      string `json:"volume"`
		FileName    string `json:"fileName"`
		Locale      string `json:"locale"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		Version     int    `json:"version"`
	} `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

type CoverListResponse struct {
	Result   string     `json:"result"`
	Response string     `json:"response"`
	Data     []CoverArt `json:"data"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Total    int        `json:"total"`
}
