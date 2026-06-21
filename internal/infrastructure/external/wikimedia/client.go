package wikimedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://ja.wikipedia.org/w/api.php"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type queryResponse struct {
	Query struct {
		Pages map[string]struct {
			Thumbnail *struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
		} `json:"pages"`
	} `json:"query"`
}

func (c *Client) FetchImage(ctx context.Context, spotName string) (*string, error) {
	params := url.Values{
		"action":      {"query"},
		"titles":      {spotName},
		"prop":        {"pageimages"},
		"format":      {"json"},
		"pithumbsize": {"800"},
		"redirects":   {"1"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("wikimedia: create request: %w", err)
	}
	req.Header.Set("User-Agent", "date-courses-go/1.0 (https://github.com/daisuke-harada/date-courses-go)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikimedia: do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wikimedia: read body: %w", err)
	}

	var result queryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("wikimedia: unmarshal response: %w", err)
	}

	for _, page := range result.Query.Pages {
		if page.Thumbnail != nil {
			return &page.Thumbnail.Source, nil
		}
	}
	return nil, nil
}
