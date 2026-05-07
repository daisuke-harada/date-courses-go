package repository

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

// DateSpotSearchParams はdate_spotsの検索条件を表します。
type DateSpotSearchParams struct {
	Name         *string
	PrefectureID *int
	GenreID      *int
	ComeTime     *string
}

// DateSpotLinkCheckParams はリンクチェック対象の検索条件を表します。
type DateSpotLinkCheckParams struct {
	// Source が "ai" のスポットのみ対象
	Source string
	// link_status が "unchecked" または LastCheckedAt が指定日数より古いもの
	DaysThreshold int
}

type DateSpotRepository interface {
	Create(ctx context.Context, dateSpot *model.DateSpot) error
	FindByID(ctx context.Context, id uint) (*model.DateSpot, error)
	Search(ctx context.Context, params DateSpotSearchParams) ([]*model.DateSpot, error)
	Update(ctx context.Context, id uint, dateSpot *model.DateSpot) error
	Delete(ctx context.Context, id uint) error

	// ExistsByNormalizedNameAndPrefecture は normalized_name + prefecture_id の組み合わせで
	// 既存レコードが存在するかどうかを確認します（AI スポットの重複防止用）。
	ExistsByNormalizedNameAndPrefecture(ctx context.Context, normalizedName string, prefectureID int) (bool, error)

	// CountByPrefectureAndGenre は指定した都道府県 × ジャンルのスポット数を返します（source 問わず）。
	CountByPrefectureAndGenre(ctx context.Context, prefectureID int, genreID int) (int64, error)

	// FindForLinkCheck はリンクチェック対象（source='ai' かつ unchecked または古いもの）を返します。
	FindForLinkCheck(ctx context.Context, params DateSpotLinkCheckParams) ([]*model.DateSpot, error)

	// UpdateLinkStatus は指定した ID の link_status と last_checked_at を更新します。
	UpdateLinkStatus(ctx context.Context, id uint, linkStatus string) error
}
