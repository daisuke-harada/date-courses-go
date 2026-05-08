package external

import (
	"context"
	"fmt"

	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/gemini"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/google_places"
	"github.com/daisuke-harada/date-courses-go/internal/infrastructure/external/nominatim"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
)

// SpotFetcherImpl は usecase.SpotFetcher の実装です。
// Gemini / Google Places / Nominatim を組み合わせてスポット情報を収集します。
type SpotFetcherImpl struct {
	gemini  *gemini.Client
	places  *google_places.Client
	nomina  *nominatim.Client
}

func NewSpotFetcher(geminiClient *gemini.Client, placesClient *google_places.Client) *SpotFetcherImpl {
	return &SpotFetcherImpl{
		gemini: geminiClient,
		places: placesClient,
		nomina: nominatim.NewClient(),
	}
}

func (f *SpotFetcherImpl) FetchSpotCandidates(ctx context.Context, prefectureName, genreName string, count int) ([]usecase.SpotCandidate, error) {
	prompt := gemini.BuildDateSpotsPrompt(prefectureName, genreName, count)

	text, err := f.gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("spot_fetcher: generate content: %w", err)
	}

	parsed, err := gemini.ParseDateSpotCandidates(text)
	if err != nil {
		return nil, fmt.Errorf("spot_fetcher: parse candidates: %w", err)
	}

	candidates := make([]usecase.SpotCandidate, 0, len(parsed))
	for _, p := range parsed {
		candidates = append(candidates, usecase.SpotCandidate{
			Name:     p.Name,
			CityName: p.CityName,
		})
	}
	return candidates, nil
}

func (f *SpotFetcherImpl) FetchPlaceDetail(ctx context.Context, spotName, cityName string) (*usecase.PlaceDetail, error) {
	d, err := f.places.FetchPlaceDetail(ctx, spotName, cityName)
	if err != nil {
		return nil, err
	}
	return &usecase.PlaceDetail{
		PhotoURL: d.PhotoURL,
	}, nil
}

func (f *SpotFetcherImpl) FetchCoordinate(ctx context.Context, spotName, cityName string) (*usecase.Coordinate, error) {
	coord, err := f.nomina.Search(ctx, spotName, cityName)
	if err != nil {
		return nil, err
	}
	if coord == nil {
		return nil, nil
	}
	return &usecase.Coordinate{Lat: coord.Lat, Lon: coord.Lon}, nil
}
