package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sangnt1552314/mangadex-tui/internal/models"
)

const (
	baseURL = "https://api.mangadex.org"
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

func GetManga(limit int, mangaType string) ([]models.Manga, error) {
	client := NewClient()

	url := getMangaDexApiUrl(limit, mangaType)

	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

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

	var mangas []models.Manga
	for _, item := range mangaList.Data {
		manga := models.Manga{
			ID:            item.ID,
			Title:         item.Attributes.Title["en"],
			Description:   item.Attributes.Description["en"],
			Status:        item.Attributes.Status,
			Year:          item.Attributes.Year,
			Tags:          make([]models.Tag, len(item.Attributes.Tags)),
			Relationships: item.Relationships,
		}
		mangas = append(mangas, manga)
	}

	if len(mangas) == 0 {
		return nil, fmt.Errorf("no manga found")
	}

	return mangas, nil
}

func getMangaDexApiUrl(limit int, mangaType string) string {
	switch mangaType {
	case "popular":
		return fmt.Sprintf("/manga?limit=%d&order[rating]=desc", limit)
	case "latest":
		return fmt.Sprintf("/manga?limit=%d&order[createdAt]=desc", limit)
	case "feature":
		return fmt.Sprintf("/manga?limit=%d&order[followedCount]=desc", limit)
	default:
		return fmt.Sprintf("/manga?limit=%d&order[rating]=desc", limit)
	}
}
