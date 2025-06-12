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

func GetMangaCover(mangaID string) (string, error) {
	client := NewClient()

	url := fmt.Sprintf("/manga/%s/cover", mangaID)
}

// getMangaApiUrl constructs the API URL for fetching manga based on the provided parameters.
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
