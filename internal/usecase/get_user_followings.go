package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// GetUserFollowingsInputPort はユーザーのフォロー一覧取得ユースケースの入力ポートです。
type GetUserFollowingsInputPort interface {
	Execute(context.Context, GetUserFollowingsInput) (*GetUserFollowingsOutput, error)
}

type GetUserFollowingsInput struct {
	UserID uint
}

type GetUserFollowingsOutput struct {
	Users []*UserWithRelations
}

type GetUserFollowingsInteractor struct {
	UserRepository           repository.UserRepository
	RelationshipRepository   repository.RelationshipRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetUserFollowingsUsecase(
	userRepository repository.UserRepository,
	relationshipRepository repository.RelationshipRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetUserFollowingsInputPort {
	return &GetUserFollowingsInteractor{
		UserRepository:           userRepository,
		RelationshipRepository:   relationshipRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetUserFollowingsInteractor) Execute(ctx context.Context, input GetUserFollowingsInput) (*GetUserFollowingsOutput, error) {
	// ユーザーの存在確認
	if _, err := i.UserRepository.FindByID(ctx, input.UserID); err != nil {
		return nil, apperror.NotFound()
	}

	followings, err := i.RelationshipRepository.FindFollowingsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result := make([]*UserWithRelations, 0, len(followings))
	for _, user := range followings {
		uwr, err := buildUserWithRelations(ctx, i.UserRepository, i.CourseRepository, i.DateSpotReviewRepository, user)
		if err != nil {
			return nil, err
		}
		result = append(result, uwr)
	}

	return &GetUserFollowingsOutput{Users: result}, nil
}
