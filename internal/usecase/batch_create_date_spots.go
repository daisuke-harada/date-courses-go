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

// SpotCandidate は外部 API から取得した実在スポット候補です。
type SpotCandidate struct {
	Name      string
	CityName  string
	Latitude  *float64
	Longitude *float64
	ImageURL  *string
	PageURL   string
}

// SpotFetcher は外部 API からスポット情報を取得するインターフェースです。
type SpotFetcher interface {
	FetchSpots(ctx context.Context, prefCode string, prefectureName string, genreID int, count int) ([]SpotCandidate, error)
}

// BatchCreateDateSpotsInput はバッチ実行の入力パラメータです。
type BatchCreateDateSpotsInput struct {
	PrefectureID   int
	PrefectureName string
	PrefCode       string
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

	candidates, err := i.fetcher.FetchSpots(ctx, input.PrefCode, input.PrefectureName, input.GenreID, i.spotsPerRun)
	if err != nil {
		return fmt.Errorf("batch: fetch spots prefecture=%s genre=%s: %w",
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
	var source model.DateSpotSource
	var mapsURL *string

	if c.PageURL != "" {
		mapsURL = &c.PageURL
	} else {
		u := BuildMapsURL(c.Name, input.PrefectureName)
		mapsURL = &u
	}

	// ジャンルIDによってソースを判定（hotpepper: 2,3,7,8,9 / jalan: それ以外）
	switch input.GenreID {
	case 2, 3, 7, 8, 9:
		source = model.DateSpotSourceHotPepper
	default:
		source = model.DateSpotSourceJalan
	}

	spot := &model.DateSpot{
		Name:           c.Name,
		CityName:       c.CityName,
		GenreID:        &input.GenreID,
		PrefectureID:   &input.PrefectureID,
		Source:         source,
		MapsURL:        mapsURL,
		NormalizedName: normalized,
		Image:          c.ImageURL,
		Latitude:       c.Latitude,
		Longitude:      c.Longitude,
	}
	return spot
}

// NormalizeName は重複チェック用に名前を正規化します（全角→半角・スペース除去・小文字化）。
func NormalizeName(name string) string {
	normalized := norm.NFKC.String(name)
	normalized = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, normalized)
	return strings.ToLower(normalized)
}

// BuildMapsURL は Google Maps 検索 URL を組み立てます。
func BuildMapsURL(spotName, prefectureName string) string {
	query := url.QueryEscape(spotName + " " + prefectureName)
	return "https://www.google.com/maps/search/" + query
}
