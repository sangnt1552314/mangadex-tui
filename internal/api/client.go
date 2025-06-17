package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sangnt1552314/mangadex-tui/internal/models"
)

const (
	baseURL           = "https://api.mangadex.org"
	coverBaseURL      = "https://uploads.mangadex.org/covers"
	ThumbnailOriginal = 0   // Original size
	Thumbnail256      = 256 // 256px width
	Thumbnail512      = 512 // 512px width
)

type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: status code %d", e.StatusCode)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: resp.StatusCode}
	}

	return resp, nil
}

func (c *Client) Post(url string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: resp.StatusCode}
	}

	return resp, nil
}

func GetManga(params models.MangaQueryParams) ([]models.Manga, error) {
	client := NewClient()

	url := getMangaApiUrl(params)

	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var mangaList models.MangaListResponse
	if err := json.NewDecoder(resp.Body).Decode(&mangaList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if mangaList.Result != "ok" {
		return nil, fmt.Errorf("API error: %s", mangaList.Result)
	}

	if len(mangaList.Data) == 0 {
		return nil, fmt.Errorf("no manga found")
	}

	return mangaList.Data, nil
}

func GetChapters(params models.ChapterQueryParams) ([]models.Chapter, error) {
	client := NewClient()

	url := getChapterApiUrl(params)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chapterList models.ChapterListResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapterList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if chapterList.Result != "ok" {
		return nil, fmt.Errorf("API error: %s", chapterList.Result)
	}

	if len(chapterList.Data) == 0 {
		return nil, fmt.Errorf("no chapters found")
	}

	return chapterList.Data, nil
}

func GetChapterListResponse(params models.ChapterQueryParams) (*models.ChapterListResponse, error) {
	client := NewClient()

	url := getChapterApiUrl(params)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chapterList models.ChapterListResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapterList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if chapterList.Result != "ok" {
		return nil, fmt.Errorf("API error: %s", chapterList.Result)
	}

	return &chapterList, nil
}

func GetMangaCover(mangaID string) (*models.CoverListResponse, error) {
	client := NewClient()

	url := fmt.Sprintf("/cover?manga[]=%s", mangaID)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coverList models.CoverListResponse
	if err := json.NewDecoder(resp.Body).Decode(&coverList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if coverList.Result != "ok" {
		return nil, fmt.Errorf("API error: %s", coverList.Result)
	}

	return &coverList, nil
}

func GetCoverURL(mangaID string, filename string, size int) string {
	url := fmt.Sprintf("%s/%s/%s", coverBaseURL, mangaID, filename)

	switch size {
	case Thumbnail256:
		return url + ".256.jpg"
	case Thumbnail512:
		return url + ".512.jpg"
	default:
		return url
	}
}

func getMangaApiUrl(params models.MangaQueryParams) string {
	queryParams := fmt.Sprintf("?limit=%d", params.Limit)

	for orderKey, orderDir := range params.Order {
		queryParams += fmt.Sprintf("&order[%s]=%s", orderKey, orderDir)
	}

	for _, rating := range params.ContentRating {
		queryParams += fmt.Sprintf("&contentRating[]=%s", rating)
	}

	for _, include := range params.Includes {
		queryParams += fmt.Sprintf("&includes[]=%s", include)
	}

	if params.HasChapters {
		queryParams += "&hasAvailableChapters=true"
	}

	url := "/manga" + queryParams

	return url
}

func getChapterApiUrl(params models.ChapterQueryParams) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", params.Limit, params.Offset)

	if len(params.Ids) > 0 {
		for _, id := range params.Ids {
			queryParams += fmt.Sprintf("&ids[]=%s", id)
		}
	}

	if params.MangaId != "" {
		queryParams += fmt.Sprintf("&manga=%s", params.MangaId)
	}

	for _, lang := range params.TranslatedLanguage {
		queryParams += fmt.Sprintf("&translatedLanguage[]=%s", lang)
	}

	url := "/chapter" + queryParams

	return url
}
