package usecase

import (
	"context"
	"sync"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
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
	UserWithRelations *model.UserWithRelations
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

	uwr, err := i.buildUserWithRelations(ctx, user)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{UserWithRelations: uwr}, nil
}

func (i *GetUserInteractor) buildUserWithRelations(ctx context.Context, user *model.User) (*model.UserWithRelations, error) {
	var (
		followerIDs  []int
		followingIDs []int
		courses      []*model.Course
		reviews      []*model.DateSpotReview
	)

	errCh := make(chan error, 4)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		followerIDs, err = i.UserRepository.FindFollowerIDsByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		followingIDs, err = i.UserRepository.FindFollowingIDsByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		courses, err = i.CourseRepository.FindByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		reviews, err = i.DateSpotReviewRepository.FindByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	return &model.UserWithRelations{
		User:         user,
		FollowerIDs:  followerIDs,
		FollowingIDs: followingIDs,
		Courses:      courses,
		Reviews:      reviews,
	}, nil
}
