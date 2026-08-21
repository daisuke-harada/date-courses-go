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

// AppGenreToHotPepper はアプリのジャンルID → HotPepper ジャンルコードのマッピングです。
// アプリのジャンル名と HotPepper の実カテゴリが一致するよう 1:1 で対応させています。
// （コード↔名称は HotPepper ジャンルマスタAPI に準拠）
var AppGenreToHotPepper = map[int]GenreCode{
	1:  "G001", // 居酒屋
	2:  "G002", // ダイニングバー・バル
	3:  "G014", // カフェ・スイーツ
	4:  "G004", // 和食
	5:  "G005", // 洋食
	6:  "G006", // イタリアン・フレンチ
	7:  "G007", // 中華
	8:  "G008", // 焼肉・ホルモン
	9:  "G013", // ラーメン
	10: "G009", // アジア・エスニック料理
	11: "G017", // 韓国料理
	12: "G012", // バー・カクテル
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
			Name    string  `json:"name"`
			Address string  `json:"address"`
			Lat     float64 `json:"lat"`
			Lng     float64 `json:"lng"`
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

// Search は HotPepper グルメ API を検索します。
// HotPepper には pref_code パラメータが無いため、地域の絞り込みは keyword（都道府県名）で行います。
func (c *Client) Search(ctx context.Context, prefectureName string, genreID int, count int) ([]Spot, error) {
	params := url.Values{
		"key":     {c.apiKey},
		"keyword": {prefectureName},
		"count":   {fmt.Sprintf("%d", count)},
		"format":  {"json"},
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
			Name:     s.Name,
			CityName: s.Address,
			Lat:      s.Lat,
			Lng:      s.Lng,
			PageURL:  s.URLs.PC,
		}
		if s.Photo.PC.L != "" {
			img := s.Photo.PC.L
			spot.ImageURL = &img
		}
		spots = append(spots, spot)
	}
	return spots, nil
}
