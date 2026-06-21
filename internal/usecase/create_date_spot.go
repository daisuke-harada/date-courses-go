package usecase

import (
	"context"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// CreateDateSpotInputPort はデートスポット作成ユースケースの入力ポートです。
type CreateDateSpotInputPort interface {
	Execute(context.Context, CreateDateSpotInput) (*CreateDateSpotOutput, error)
}

// CreateDateSpotInput はデートスポッ���作成の入力データです。
type CreateDateSpotInput struct {
	Name         string
	GenreID      int
	PrefectureID int
	CityName     string
	Image        *string
}

// Validate はデートスポット��成の入力デー��をバリデーションします��
func (i *CreateDateSpotInput) Validate() error {
	var errs []string

	if strings.TrimSpace(i.Name) == "" {
		errs = append(errs, "スポット名を入力してください")
	}
	if i.GenreID <= 0 {
		errs = append(errs, "���ャンルを選択してくだ���い")
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
	if err := input.Validate(); err != nil {
		return nil, err
	}

	dateSpot := &model.DateSpot{
		Name:         input.Name,
		GenreID:      &input.GenreID,
		PrefectureID: &input.PrefectureID,
		CityName:     input.CityName,
		Image:        input.Image,
	}

	if err := i.DateSpotRepository.Create(ctx, dateSpot); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &CreateDateSpotOutput{DateSpotID: dateSpot.ID}, nil
}
