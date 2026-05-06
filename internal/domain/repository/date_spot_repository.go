package repository

import (
	"context"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DateSpotSearchParams はdate_spotsの検索条件を表します。
type DateSpotSearchParams struct {
	Name         *string
	PrefectureID *int
	GenreID      *int
	ComeTime     *string
}

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *model.DateSpot) error
	FindByID(ctx context.Context, id uint) (*model.DateSpot, error)
	Search(ctx context.Context, params DateSpotSearchParams) ([]*model.DateSpot, error)
	Update(ctx context.Context, id uint, dateSpot *model.DateSpot) error
	Delete(ctx context.Context, id uint) error

	// AI バッチ統合用メソッド

	// CountByPrefectureAndGenre は都道府県ID×ジャンルIDのスポット総数（source 問わず）を返します。
	CountByPrefectureAndGenre(ctx context.Context, prefectureID, genreID int) (int64, error)
	// ExistsByNormalizedNameAndPrefecture は正規化名と都道府県IDでスポットの存在チェックをします（AI スポットのみ）。
	ExistsByNormalizedNameAndPrefecture(ctx context.Context, normalizedName string, prefectureID int) (bool, error)
	// FindAISpotsForLinkCheck はリンクチェックが必要な AI スポット（unchecked or 7日以上未チェック）を返します。
	FindAISpotsForLinkCheck(ctx context.Context) ([]*model.DateSpot, error)
	// UpdateLinkStatus は AI スポットのリンクステータスと最終確認日時を更新します。
	UpdateLinkStatus(ctx context.Context, id uint, linkStatus string, lastCheckedAt time.Time) error
}
