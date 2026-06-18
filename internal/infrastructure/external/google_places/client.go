package google_places

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const searchURL = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
const detailURL = "https://maps.googleapis.com/maps/api/place/details/json"
const photoURL = "https://maps.googleapis.com/maps/api/place/photo"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PlaceDetail は取得したスポットの詳細情報です。
type PlaceDetail struct {
	PhotoURL    *string
	OpeningTime *time.Time
	ClosingTime *time.Time
}

type findPlaceResponse struct {
	Candidates []struct {
		PlaceID string `json:"place_id"`
	} `json:"candidates"`
	Status string `json:"status"`
}

type placeDetailResponse struct {
	Result struct {
		Photos []struct {
			PhotoReference string `json:"photo_reference"`
		} `json:"photos"`
		OpeningHours *struct {
			Periods []struct {
				Open  *periodTime `json:"open"`
				Close *periodTime `json:"close"`
			} `json:"periods"`
		} `json:"opening_hours"`
	} `json:"result"`
	Status string `json:"status"`
}

type periodTime struct {
	Time string `json:"time"` // "HHMM" 形式
}

func (c *Client) FetchPlaceDetail(ctx context.Context, spotName, cityName string) (*PlaceDetail, error) {
	placeID, err := c.findPlaceID(ctx, spotName+" "+cityName)
	if err != nil {
		return nil, err
	}
	if placeID == "" {
		return &PlaceDetail{}, nil
	}

	return c.fetchDetail(ctx, placeID)
}

func (c *Client) findPlaceID(ctx context.Context, query string) (string, error) {
	params := url.Values{
		"input":          {query},
		"inputtype":      {"textquery"},
		"fields":         {"place_id"},
		"language":       {"ja"},
		"key":            {c.apiKey},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("google_places: create find request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("google_places: do find request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("google_places: read find body: %w", err)
	}

	var result findPlaceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("google_places: unmarshal find response: %w", err)
	}

	if result.Status != "OK" || len(result.Candidates) == 0 {
		slog.InfoContext(ctx, "google_places: no place found", "query", query, "status", result.Status)
		return "", nil
	}

	return result.Candidates[0].PlaceID, nil
}

func (c *Client) fetchDetail(ctx context.Context, placeID string) (*PlaceDetail, error) {
	params := url.Values{
		"place_id": {placeID},
		"fields":   {"photos,opening_hours"},
		"language": {"ja"},
		"key":      {c.apiKey},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("google_places: create detail request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google_places: do detail request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google_places: read detail body: %w", err)
	}

	var result placeDetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("google_places: unmarshal detail response: %w", err)
	}

	detail := &PlaceDetail{}

	if len(result.Result.Photos) > 0 {
		photoRef := result.Result.Photos[0].PhotoReference
		photoURLStr := fmt.Sprintf("%s?maxwidth=800&photoreference=%s&key=%s", photoURL, photoRef, c.apiKey)
		detail.PhotoURL = &photoURLStr
	}

	if result.Result.OpeningHours != nil && len(result.Result.OpeningHours.Periods) > 0 {
		p := result.Result.OpeningHours.Periods[0]
		if p.Open != nil {
			t := parseHHMM(p.Open.Time)
			detail.OpeningTime = &t
		}
		if p.Close != nil {
			t := parseHHMM(p.Close.Time)
			detail.ClosingTime = &t
		}
	}

	return detail, nil
}

// parseHHMM は "0900" 形式の文字列を time.Time に変換します（日付は今日の UTC）。
func parseHHMM(hhmm string) time.Time {
	if len(hhmm) != 4 {
		return time.Time{}
	}
	t, _ := time.Parse("1504", hhmm)
	return t
}
