package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const searchURL = "https://nominatim.openstreetmap.org/search"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Coordinate struct {
	Lat float64
	Lon float64
}

type nominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// Search はスポット名と都市名から緯度経度を返します。
// 見つからない場合は nil を返します（エラーではない）。
func (c *Client) Search(ctx context.Context, spotName, cityName string) (*Coordinate, error) {
	params := url.Values{
		"q":              {spotName + " " + cityName},
		"format":         {"json"},
		"limit":          {"1"},
		"accept-language": {"ja"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("nominatim: create request: %w", err)
	}
	// Nominatim の利用規約: User-Agent を設定すること
	req.Header.Set("User-Agent", "date-courses-go/1.0 (https://github.com/daisuke-harada/date-courses-go)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nominatim: do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nominatim: read body: %w", err)
	}

	var results []nominatimResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("nominatim: unmarshal response: %w", err)
	}

	if len(results) == 0 {
		slog.InfoContext(ctx, "nominatim: no result", "spot", spotName, "city", cityName)
		return nil, nil
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("nominatim: parse lat: %w", err)
	}
	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("nominatim: parse lon: %w", err)
	}

	return &Coordinate{Lat: lat, Lon: lon}, nil
}
