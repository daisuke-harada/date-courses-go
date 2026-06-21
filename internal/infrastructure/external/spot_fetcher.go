package external

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/hotpepper"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/jalan"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/wikimedia"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
)

// SpotFetcherImpl は usecase.SpotFetcher の実装です。
// ジャンルに応じて HotPepper / じゃらん を使い分け、画像は Wikimedia でフォールバックします。
type SpotFetcherImpl struct {
	hotpepper *hotpepper.Client
	jalan     *jalan.Client
	wikimedia *wikimedia.Client
}

func NewSpotFetcher(hp *hotpepper.Client, jl *jalan.Client, wm *wikimedia.Client) *SpotFetcherImpl {
	return &SpotFetcherImpl{
		hotpepper: hp,
		jalan:     jl,
		wikimedia: wm,
	}
}

func (f *SpotFetcherImpl) FetchSpots(ctx context.Context, prefCode string, prefectureName string, genreID int, count int) ([]usecase.SpotCandidate, error) {
	var spots []usecase.SpotCandidate
	var err error

	if hotpepper.IsHotPepperGenre(genreID) {
		spots, err = f.fetchFromHotPepper(ctx, prefCode, genreID, count)
	} else if jalan.IsJalanGenre(genreID) {
		spots, err = f.fetchFromJalan(ctx, prefCode, genreID, count)
	} else {
		return nil, fmt.Errorf("spot_fetcher: unsupported genre id: %d", genreID)
	}

	if err != nil {
		return nil, err
	}

	// 画像がないスポットは Wikimedia でフォールバック
	for i := range spots {
		if spots[i].ImageURL == nil {
			imageURL, wErr := f.wikimedia.FetchImage(ctx, spots[i].Name)
			if wErr != nil {
				slog.InfoContext(ctx, "spot_fetcher: wikimedia fallback failed", "name", spots[i].Name, "err", wErr)
				continue
			}
			spots[i].ImageURL = imageURL
		}
	}

	return spots, nil
}

func (f *SpotFetcherImpl) fetchFromHotPepper(ctx context.Context, prefCode string, genreID int, count int) ([]usecase.SpotCandidate, error) {
	results, err := f.hotpepper.Search(ctx, prefCode, genreID, count)
	if err != nil {
		return nil, fmt.Errorf("spot_fetcher: hotpepper search: %w", err)
	}

	candidates := make([]usecase.SpotCandidate, 0, len(results))
	for _, s := range results {
		c := usecase.SpotCandidate{
			Name:     s.Name,
			CityName: s.CityName,
			PageURL:  s.PageURL,
			ImageURL: s.ImageURL,
		}
		if s.Lat != 0 {
			lat := s.Lat
			c.Latitude = &lat
		}
		if s.Lng != 0 {
			lng := s.Lng
			c.Longitude = &lng
		}
		candidates = append(candidates, c)
	}
	return candidates, nil
}

func (f *SpotFetcherImpl) fetchFromJalan(ctx context.Context, prefCode string, genreID int, count int) ([]usecase.SpotCandidate, error) {
	results, err := f.jalan.Search(ctx, prefCode, genreID, count)
	if err != nil {
		return nil, fmt.Errorf("spot_fetcher: jalan search: %w", err)
	}

	candidates := make([]usecase.SpotCandidate, 0, len(results))
	for _, s := range results {
		c := usecase.SpotCandidate{
			Name:     s.Name,
			CityName: s.CityName,
			PageURL:  s.PageURL,
			ImageURL: s.ImageURL,
		}
		if s.Lat != 0 {
			lat := s.Lat
			c.Latitude = &lat
		}
		if s.Lng != 0 {
			lng := s.Lng
			c.Longitude = &lng
		}
		candidates = append(candidates, c)
	}
	return candidates, nil
}
