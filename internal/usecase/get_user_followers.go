package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// GetUserFollowersInputPort はユーザーのフォロワー一覧取得ユースケースの入力ポートです。
type GetUserFollowersInputPort interface {
	Execute(context.Context, GetUserFollowersInput) (*GetUserFollowersOutput, error)
}

type GetUserFollowersInput struct {
	UserID uint
}

type GetUserFollowersOutput struct {
	Users []*UserWithRelations
}

type GetUserFollowersInteractor struct {
	UserRepository           repository.UserRepository
	RelationshipRepository   repository.RelationshipRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetUserFollowersUsecase(
	userRepository repository.UserRepository,
	relationshipRepository repository.RelationshipRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetUserFollowersInputPort {
	return &GetUserFollowersInteractor{
		UserRepository:           userRepository,
		RelationshipRepository:   relationshipRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetUserFollowersInteractor) Execute(ctx context.Context, input GetUserFollowersInput) (*GetUserFollowersOutput, error) {
	// ユーザーの存在確認
	if _, err := i.UserRepository.FindByID(ctx, input.UserID); err != nil {
		return nil, apperror.NotFound()
	}

	followers, err := i.RelationshipRepository.FindFollowersByUserID(ctx, input.UserID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result := make([]*UserWithRelations, 0, len(followers))
	for _, user := range followers {
		uwr, err := buildUserWithRelations(ctx, i.UserRepository, i.CourseRepository, i.DateSpotReviewRepository, user)
		if err != nil {
			return nil, err
		}
		result = append(result, uwr)
	}

	return &GetUserFollowersOutput{Users: result}, nil
}
