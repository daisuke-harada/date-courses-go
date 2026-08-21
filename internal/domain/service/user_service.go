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
	BuildUsersWithRelations(ctx context.Context, users []*model.User) ([]*model.UserWithRelations, error)
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

// BuildUserWithRelations はユーザーに紐づく関連データを取得して UserWithRelations を返します。
func (s *userService) BuildUserWithRelations(ctx context.Context, user *model.User) (*model.UserWithRelations, error) {
	result, err := s.BuildUsersWithRelations(ctx, []*model.User{user})
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// BuildUsersWithRelations は複数ユーザーの関連データをまとめて取得します。
// 関連ごとに1回ずつ（フォロワー・フォロー中・コース・レビューの計4回）しかクエリを
// 発行しません。ユーザー1人につき4回投げると、一覧では人数に比例して
// クエリが増えてしまうためです。
func (s *userService) BuildUsersWithRelations(ctx context.Context, users []*model.User) ([]*model.UserWithRelations, error) {
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}

	var (
		followerIDs  map[uint][]int
		followingIDs map[uint][]int
		courses      map[uint][]*model.Course
		reviews      map[uint][]*model.DateSpotReview
	)

	errCh := make(chan error, 4)
	var wg sync.WaitGroup

	run := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(); err != nil {
				errCh <- apperror.InternalServerError(err)
			}
		}()
	}

	run(func() error {
		var err error
		followerIDs, err = s.UserRepository.FindFollowerIDsByUserIDs(ctx, userIDs)
		return err
	})
	run(func() error {
		var err error
		followingIDs, err = s.UserRepository.FindFollowingIDsByUserIDs(ctx, userIDs)
		return err
	})
	run(func() error {
		var err error
		courses, err = s.CourseRepository.FindPublicByUserIDs(ctx, userIDs)
		return err
	})
	run(func() error {
		var err error
		reviews, err = s.DateSpotReviewRepository.FindByUserIDs(ctx, userIDs)
		return err
	})

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	result := make([]*model.UserWithRelations, 0, len(users))
	for _, u := range users {
		result = append(result, &model.UserWithRelations{
			User:         u,
			FollowerIDs:  emptyIfNil(followerIDs[u.ID]),
			FollowingIDs: emptyIfNil(followingIDs[u.ID]),
			Courses:      courses[u.ID],
			Reviews:      reviews[u.ID],
		})
	}
	return result, nil
}

// emptyIfNil は nil スライスを空スライスに変換します。
// レスポンスで null ではなく [] を返すために必要です。
func emptyIfNil(ids []int) []int {
	if ids == nil {
		return []int{}
	}
	return ids
}
