package usecase

import (
	"context"
	"sync"

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
	Users []*model.UserWithRelations
}

type GetUsersInteractor struct {
	UserRepository           repository.UserRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewGetUsersUsecase(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) GetUsersInputPort {
	return &GetUsersInteractor{
		UserRepository:           userRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *GetUsersInteractor) Execute(ctx context.Context, input GetUsersInput) (*GetUsersOutput, error) {
	users, err := i.UserRepository.Search(ctx, input.Name)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	result := make([]*model.UserWithRelations, len(users))
	errCh := make(chan error, len(users))
	var wg sync.WaitGroup

	for idx, user := range users {
		wg.Add(1)
		go func(idx int, user *model.User) {
			defer wg.Done()
			uwr, err := i.buildUserWithRelations(ctx, user)
			if err != nil {
				errCh <- err
				return
			}
			result[idx] = uwr
		}(idx, user)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	return &GetUsersOutput{Users: result}, nil
}

func (i *GetUsersInteractor) buildUserWithRelations(ctx context.Context, user *model.User) (*model.UserWithRelations, error) {
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
