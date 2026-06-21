package usecase

import (
	"context"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// UpdateDateSpotInputPort はデートスポット更新ユースケースの入力ポートです。
type UpdateDateSpotInputPort interface {
	Execute(context.Context, UpdateDateSpotInput) error
}

// UpdateDateSpotInput はデートスポット更新の入力データです。
type UpdateDateSpotInput struct {
	DateSpotID   uint
	Name         string
	GenreID      int
	PrefectureID int
	CityName     string
	Image        *string
}

// Validate はデートスポット更新の入力データをバリデーションします。
func (i *UpdateDateSpotInput) Validate() error {
	var errs []string

	if strings.TrimSpace(i.Name) == "" {
		errs = append(errs, "スポット名を入力してください")
	}
	if i.GenreID <= 0 {
		errs = append(errs, "ジャンルを選択してください")
	}
	if i.PrefectureID <= 0 {
		errs = append(errs, "都道府県を選択してください")
	}
	if strings.TrimSpace(i.CityName) == "" {
		errs = append(errs, "市区町村を入力してください")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}
	return nil
}

type UpdateDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewUpdateDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
) UpdateDateSpotInputPort {
	return &UpdateDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *UpdateDateSpotInteractor) Execute(ctx context.Context, input UpdateDateSpotInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	dateSpot := &model.DateSpot{
		Name:         input.Name,
		GenreID:      &input.GenreID,
		PrefectureID: &input.PrefectureID,
		CityName:     input.CityName,
		Image:        input.Image,
	}

	if err := i.DateSpotRepository.Update(ctx, input.DateSpotID, dateSpot); err != nil {
		return apperror.InternalServerError(err)
	}

	return nil
}
