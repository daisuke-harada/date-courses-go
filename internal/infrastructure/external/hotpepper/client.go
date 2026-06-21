package hotpepper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"

// GenreCode は HotPepper Gourmet のジャンルコードです。
type GenreCode string

const (
	GenreIzakaya  GenreCode = "G001" // 居酒屋
	GenreCafe     GenreCode = "G014" // カフェ・スイーツ
	GenreSushi    GenreCode = "G013" // 寿司
	GenreYakiniku GenreCode = "G009" // 焼肉・ホルモン
)

// AppGenreToHotPepper はアプリのジャンルID → HotPepper ジャンルコードのマッピングです。
var AppGenreToHotPepper = map[int]GenreCode{
	2: "",           // 飲食店（ジャンル指定なし）
	3: GenreCafe,
	7: GenreSushi,
	8: GenreIzakaya,
	9: GenreYakiniku,
}

// IsHotPepperGenre はジャンルID が HotPepper 対象かどうかを返します。
func IsHotPepperGenre(genreID int) bool {
	_, ok := AppGenreToHotPepper[genreID]
	return ok
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type Spot struct {
	Name     string
	CityName string
	Lat      float64
	Lng      float64
	ImageURL *string
	PageURL  string
}

type response struct {
	Results struct {
		Shop []struct {
			Name    string `json:"name"`
			Address string `json:"address"`
			Lat     string `json:"lat"`
			Lng     string `json:"lng"`
			Photo   struct {
				PC struct {
					L string `json:"l"`
				} `json:"pc"`
			} `json:"photo"`
			URLs struct {
				PC string `json:"pc"`
			} `json:"urls"`
		} `json:"shop"`
	} `json:"results"`
}

func (c *Client) Search(ctx context.Context, prefCode string, genreID int, count int) ([]Spot, error) {
	params := url.Values{
		"key":       {c.apiKey},
		"pref_code": {prefCode},
		"count":     {fmt.Sprintf("%d", count)},
		"format":    {"json"},
	}
	if gc, ok := AppGenreToHotPepper[genreID]; ok && gc != "" {
		params.Set("genre", string(gc))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("hotpepper: create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hotpepper: do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("hotpepper: read body: %w", err)
	}

	var result response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("hotpepper: unmarshal response: %w", err)
	}

	spots := make([]Spot, 0, len(result.Results.Shop))
	for _, s := range result.Results.Shop {
		spot := Spot{
			Name:    s.Name,
			PageURL: s.URLs.PC,
		}
		if s.Photo.PC.L != "" {
			img := s.Photo.PC.L
			spot.ImageURL = &img
		}
		_, _ = fmt.Sscanf(s.Lat, "%f", &spot.Lat)
		_, _ = fmt.Sscanf(s.Lng, "%f", &spot.Lng)
		spots = append(spots, spot)
	}
	return spots, nil
}
