package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// GetUserInputPort はユーザー単体取得ユースケースの入力ポートです。
type GetUserInputPort interface {
	Execute(context.Context, GetUserInput) (*GetUserOutput, error)
}

type GetUserInput struct {
	ID uint
}

type GetUserOutput struct {
	UserWithRelations *UserWithRelations
}

type GetUserInteractor struct {
	UserRepository           repository.UserRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetUserUsecase(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetUserInputPort {
	return &GetUserInteractor{
		UserRepository:           userRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetUserInteractor) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	uwr, err := buildUserWithRelations(ctx, i.UserRepository, i.CourseRepository, i.DateSpotReviewRepository, user)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{UserWithRelations: uwr}, nil
}
