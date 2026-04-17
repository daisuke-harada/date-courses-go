package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// CreateDateSpotInputPort はデートスポット作成ユースケースの入力ポートです。
type CreateDateSpotInputPort interface {
	Execute(context.Context, CreateDateSpotInput) (*CreateDateSpotOutput, error)
}

// CreateDateSpotInput はデートスポット作成の入力データです。
type CreateDateSpotInput struct {
	Name         string
	GenreID      int
	PrefectureID int
	CityName     string
	OpeningTime  *time.Time
	ClosingTime  *time.Time
	Image        *string
}

// Validate はデートスポット作成の入力データをバリデーションします。
func (i *CreateDateSpotInput) Validate() error {
	var errs []string

	// name: presence
	if strings.TrimSpace(i.Name) == "" {
		errs = append(errs, "スポット名を入力してください")
	}

	// genre_id: presence, positive
	if i.GenreID <= 0 {
		errs = append(errs, "ジャンルを選択してください")
	}

	// prefecture_id: presence, positive
	if i.PrefectureID <= 0 {
		errs = append(errs, "都道府県を選択してください")
	}

	// city_name: presence
	if strings.TrimSpace(i.CityName) == "" {
		errs = append(errs, "市区町村を入力してください")
	}

	// opening_time: presence
	if i.OpeningTime == nil {
		errs = append(errs, "営業開始時刻を入力してください")
	}

	// closing_time: presence
	if i.ClosingTime == nil {
		errs = append(errs, "営業終了時刻を入力してください")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	return nil
}

// CreateDateSpotOutput はデートスポット作成の出力データです。
type CreateDateSpotOutput struct {
	DateSpotID uint
}

type CreateDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewCreateDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
) CreateDateSpotInputPort {
	return &CreateDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *CreateDateSpotInteractor) Execute(ctx context.Context, input CreateDateSpotInput) (*CreateDateSpotOutput, error) {
	// バリデーション
	if err := input.Validate(); err != nil {
		return nil, err
	}

	dateSpot := &model.DateSpot{
		Name:         input.Name,
		GenreID:      &input.GenreID,
		PrefectureID: &input.PrefectureID,
		CityName:     input.CityName,
		OpeningTime:  input.OpeningTime,
		ClosingTime:  input.ClosingTime,
		Image:        input.Image,
	}

	if err := i.DateSpotRepository.Create(ctx, dateSpot); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &CreateDateSpotOutput{DateSpotID: dateSpot.ID}, nil
}
