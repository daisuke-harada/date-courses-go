package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// GetUsersInputPort はユーザー一覧取得ユースケースの入力ポートです。
type GetUsersInputPort interface {
	Execute(context.Context, GetUsersInput) (*GetUsersOutput, error)
}

type GetUsersInput struct {
	Name *string
}

type GetUsersOutput struct {
	Users []*UserWithRelations
}

// UserWithRelations はユーザーと関連データをまとめた中間型です。
type UserWithRelations struct {
	User         *model.User
	FollowerIDs  []int
	FollowingIDs []int
	Courses      []*model.Course
	Reviews      []*model.DateSpotReview
}

type GetUsersInteractor struct {
	UserRepository          repository.UserRepository
	CourseRepository        repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetUsersUsecase(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetUsersInputPort {
	return &GetUsersInteractor{
		UserRepository:          userRepository,
		CourseRepository:        courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetUsersInteractor) Execute(ctx context.Context, input GetUsersInput) (*GetUsersOutput, error) {
	users, err := i.UserRepository.Search(ctx, input.Name)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result := make([]*UserWithRelations, 0, len(users))
	for _, user := range users {
		uwr, err := i.buildUserWithRelations(ctx, user)
		if err != nil {
			return nil, err
		}
		result = append(result, uwr)
	}

	return &GetUsersOutput{Users: result}, nil
}

func (i *GetUsersInteractor) buildUserWithRelations(ctx context.Context, user *model.User) (*UserWithRelations, error) {
	followerIDs, err := i.UserRepository.FindFollowerIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	followingIDs, err := i.UserRepository.FindFollowingIDsByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	courses, err := i.CourseRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	reviews, err := i.DateSpotReviewRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &UserWithRelations{
		User:         user,
		FollowerIDs:  followerIDs,
		FollowingIDs: followingIDs,
		Courses:      courses,
		Reviews:      reviews,
	}, nil
}
