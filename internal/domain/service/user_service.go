package service

import (
	"context"
	"sync"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
)

// UserService はユーザーに関するドメインサービスです。
type UserService interface {
	BuildUserWithRelations(ctx context.Context, user *model.User) (*model.UserWithRelations, error)
}

type userService struct {
	UserRepository           repository.UserRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) UserService {
	return &userService{
		UserRepository:           userRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

// BuildUserWithRelations はユーザーに紐づく関連データを並列取得して UserWithRelations を返します。
func (s *userService) BuildUserWithRelations(ctx context.Context, user *model.User) (*model.UserWithRelations, error) {
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
		followerIDs, err = s.UserRepository.FindFollowerIDsByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		followingIDs, err = s.UserRepository.FindFollowingIDsByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		courses, err = s.CourseRepository.FindByUserID(ctx, user.ID)
		if err != nil {
			errCh <- apperror.InternalServerError(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		reviews, err = s.DateSpotReviewRepository.FindByUserID(ctx, user.ID)
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
