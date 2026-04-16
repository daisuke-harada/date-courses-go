package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// DeleteDateSpotInputPort はデートスポット削除ユースケースの入力ポートです。
type DeleteDateSpotInputPort interface {
	Execute(context.Context, DeleteDateSpotInput) error
}

// DeleteDateSpotInput はデートスポット削除の入力データです。
type DeleteDateSpotInput struct {
	DateSpotID uint
}

type DeleteDateSpotInteractor struct {
	DateSpotRepository repository.DateSpotRepository
}

func NewDeleteDateSpotUsecase(
	dateSpotRepository repository.DateSpotRepository,
) DeleteDateSpotInputPort {
	return &DeleteDateSpotInteractor{
		DateSpotRepository: dateSpotRepository,
	}
}

func (i *DeleteDateSpotInteractor) Execute(ctx context.Context, input DeleteDateSpotInput) error {
	if err := i.DateSpotRepository.Delete(ctx, input.DateSpotID); err != nil {
		return apperror.InternalServerError(err)
	}
	return nil
}
