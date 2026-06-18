package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"unicode"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"golang.org/x/text/unicode/norm"
)

// SpotCandidate は外部 API から取得したスポット候補です。
type SpotCandidate struct {
	Name     string
	CityName string
}

// PlaceDetail は Google Places API から取得したスポット詳細です。
type PlaceDetail struct {
	PhotoURL    *string
	OpeningTime *interface{}
	ClosingTime *interface{}
}

// Coordinate は Nominatim から取得した緯度経度です。
type Coordinate struct {
	Lat float64
	Lon float64
}

// SpotFetcher は外部 API からスポット情報を取得するインターフェースです。
type SpotFetcher interface {
	FetchSpotCandidates(ctx context.Context, prefectureName, genreName string, count int) ([]SpotCandidate, error)
	FetchPlaceDetail(ctx context.Context, spotName, cityName string) (*PlaceDetail, error)
	FetchCoordinate(ctx context.Context, spotName, cityName string) (*Coordinate, error)
}

// BatchCreateDateSpotsInput はバッチ実行の入力パラメータです。
type BatchCreateDateSpotsInput struct {
	PrefectureID   int
	PrefectureName string
	GenreID        int
	GenreName      string
}

// BatchCreateDateSpotsInteractor はデートスポット自動収集バッチのユースケースです。
type BatchCreateDateSpotsInteractor struct {
	repo             repository.DateSpotRepository
	fetcher          SpotFetcher
	minExistingSpots int
	spotsPerRun      int
}

func NewBatchCreateDateSpotsInteractor(
	repo repository.DateSpotRepository,
	fetcher SpotFetcher,
	minExistingSpots int,
	spotsPerRun int,
) *BatchCreateDateSpotsInteractor {
	return &BatchCreateDateSpotsInteractor{
		repo:             repo,
		fetcher:          fetcher,
		minExistingSpots: minExistingSpots,
		spotsPerRun:      spotsPerRun,
	}
}

// Execute は1つの都道府県×ジャンルの組み合わせに対してスポット収集を実行します。
func (i *BatchCreateDateSpotsInteractor) Execute(ctx context.Context, input BatchCreateDateSpotsInput) error {
	count, err := i.repo.CountByPrefectureAndGenre(ctx, input.PrefectureID, input.GenreID)
	if err != nil {
		return fmt.Errorf("batch: count spots: %w", err)
	}
	if count >= int64(i.minExistingSpots) {
		slog.InfoContext(ctx, "batch: skip combination, enough spots exist",
			"prefecture_id", input.PrefectureID,
			"genre_id", input.GenreID,
			"count", count,
		)
		return nil
	}

	candidates, err := i.fetcher.FetchSpotCandidates(ctx, input.PrefectureName, input.GenreName, i.spotsPerRun)
	if err != nil {
		return fmt.Errorf("batch: fetch candidates prefecture=%s genre=%s: %w",
			input.PrefectureName, input.GenreName, err)
	}

	var newSpots []*model.DateSpot
	for _, c := range candidates {
		normalized := NormalizeName(c.Name)

		exists, err := i.repo.ExistsByNormalizedNameAndPrefecture(ctx, normalized, input.PrefectureID)
		if err != nil {
			slog.ErrorContext(ctx, "batch: check existence failed, skipping", "name", c.Name, "err", err)
			continue
		}
		if exists {
			slog.InfoContext(ctx, "batch: skip duplicate spot", "name", c.Name)
			continue
		}

		spot := i.buildDateSpot(c, input, normalized)
		i.enrichWithExternalAPIs(ctx, spot, c)
		newSpots = append(newSpots, spot)
	}

	if len(newSpots) == 0 {
		slog.InfoContext(ctx, "batch: no new spots to create",
			"prefecture_id", input.PrefectureID,
			"genre_id", input.GenreID,
		)
		return nil
	}

	if err := i.repo.CreateBatch(ctx, newSpots); err != nil {
		return fmt.Errorf("batch: create batch: %w", err)
	}

	slog.InfoContext(ctx, "batch: created spots",
		"prefecture_id", input.PrefectureID,
		"genre_id", input.GenreID,
		"count", len(newSpots),
	)
	return nil
}

func (i *BatchCreateDateSpotsInteractor) buildDateSpot(
	c SpotCandidate,
	input BatchCreateDateSpotsInput,
	normalized string,
) *model.DateSpot {
	mapsURL := BuildMapsURL(c.Name, input.PrefectureName)
	linkStatus := model.DateSpotLinkStatusUnchecked
	return &model.DateSpot{
		Name:           c.Name,
		CityName:       c.CityName,
		GenreID:        &input.GenreID,
		PrefectureID:   &input.PrefectureID,
		Source:         model.DateSpotSourceAI,
		MapsURL:        &mapsURL,
		LinkStatus:     linkStatus,
		NormalizedName: normalized,
	}
}

func (i *BatchCreateDateSpotsInteractor) enrichWithExternalAPIs(ctx context.Context, spot *model.DateSpot, c SpotCandidate) {
	coord, err := i.fetcher.FetchCoordinate(ctx, c.Name, c.CityName)
	if err != nil {
		slog.InfoContext(ctx, "batch: fetch coordinate failed, skipping", "name", c.Name, "err", err)
	} else if coord != nil {
		spot.Latitude = &coord.Lat
		spot.Longitude = &coord.Lon
	}

	detail, err := i.fetcher.FetchPlaceDetail(ctx, c.Name, c.CityName)
	if err != nil {
		slog.InfoContext(ctx, "batch: fetch place detail failed, skipping", "name", c.Name, "err", err)
	} else if detail != nil {
		spot.Image = detail.PhotoURL
	}
}

// NormalizeName は重複チェック用に名前を正規化します（全角→半角・スペース除去・小文字化）。
func NormalizeName(name string) string {
	// Unicode 正規化（全角英数→半角）
	normalized := norm.NFKC.String(name)
	// スペース除去
	normalized = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, normalized)
	// 小文字化
	return strings.ToLower(normalized)
}

// BuildMapsURL は Google Maps 検索 URL を組み立てます。
func BuildMapsURL(spotName, prefectureName string) string {
	query := url.QueryEscape(spotName + " " + prefectureName)
	return "https://www.google.com/maps/search/" + query
}
