package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/repository"
)

type GetDateSpotsInputPort interface {
	Execute(context.Context, GetDateSpotsInput) (*GetDateSpotsOutput, error)
}

// GetDateSpotsInput は GetApiV1DateSpotsParams に対応するusecaseの入力型です。
// usecase層をOpenAPIの生成型に依存させないよう、独自の型として定義します。
type GetDateSpotsInput struct {
	DateSpotName *string
	PrefectureID *int
	GenreID      *int
	ComeTime     *string
}

type GetDateSpotsOutput struct {
	DateSpots []*model.DateSpot
}

type GetDateSpotsInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewGetDateSpotsUsecase(
	dateSpotRepository repository.DateSpotRepository,
) GetDateSpotsInputPort {
	return &GetDateSpotsInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *GetDateSpotsInteractor) Execute(ctx context.Context, input GetDateSpotsInput) (*GetDateSpotsOutput, error) {
	params := repository.DateSpotSearchParams{
		Name:         input.DateSpotName,
		PrefectureID: input.PrefectureID,
		GenreID:      input.GenreID,
		ComeTime:     input.ComeTime,
	}
	dateSpots, err := i.DateSpotRepository.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	return &GetDateSpotsOutput{
		DateSpots: dateSpots,
	}, nil
}
