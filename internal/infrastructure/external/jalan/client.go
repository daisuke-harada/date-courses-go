package jalan

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://webservice.recruit.co.jp/jalan/spot/v1/"

// GenreCode はじゃらんnet 観光スポット API のジャンルコードです。
type GenreCode string

const (
	GenreShopping   GenreCode = "sej06" // ショッピング
	GenreOutdoor    GenreCode = "sej04" // スポーツ・アウトドア
	GenreThemePark  GenreCode = "sej01" // テーマパーク・遊園地・動物園・水族館
	GenreLandmark   GenreCode = "sej08" // 夜景・展望・ランドマーク
	GenrePark       GenreCode = "sej03" // 自然景観・公園
)

// AppGenreToJalan はアプリのジャンルID → じゃらん ジャンルコードのマッピングです。
var AppGenreToJalan = map[int]GenreCode{
	1:  GenreShopping,
	4:  GenreOutdoor,
	5:  GenreThemePark,
	6:  GenreThemePark,
	10: GenreOutdoor,
	11: GenreLandmark,
	12: GenrePark,
}

// IsJalanGenre はジャンルID がじゃらん対象かどうかを返します。
func IsJalanGenre(genreID int) bool {
	_, ok := AppGenreToJalan[genreID]
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

type xmlResponse struct {
	XMLName xml.Name   `xml:"Results"`
	Spots   []xmlSpot  `xml:"spot"`
}

type xmlSpot struct {
	Name      string  `xml:"posName"`
	Address   string  `xml:"address"`
	Lat       float64 `xml:"latitude"`
	Lng       float64 `xml:"longitude"`
	ImagePath string  `xml:"imagePath"`
	SpotURL   string  `xml:"spotURL"`
}

func (c *Client) Search(ctx context.Context, prefCode string, genreID int, count int) ([]Spot, error) {
	gc, ok := AppGenreToJalan[genreID]
	if !ok {
		return nil, fmt.Errorf("jalan: unsupported genre id: %d", genreID)
	}

	params := url.Values{
		"key":      {c.apiKey},
		"prefCode": {prefCode},
		"genre":    {string(gc)},
		"count":    {fmt.Sprintf("%d", count)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("jalan: create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jalan: do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("jalan: read body: %w", err)
	}

	var result xmlResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("jalan: unmarshal response: %w", err)
	}

	spots := make([]Spot, 0, len(result.Spots))
	for _, s := range result.Spots {
		spot := Spot{
			Name:    s.Name,
			Lat:     s.Lat,
			Lng:     s.Lng,
			PageURL: s.SpotURL,
		}
		if s.ImagePath != "" {
			img := s.ImagePath
			spot.ImageURL = &img
		}
		spots = append(spots, spot)
	}
	return spots, nil
}
