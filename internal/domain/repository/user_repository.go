package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	Search(ctx context.Context, name *string) ([]*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// FindFollowerIDsByUserIDs / FindFollowingIDsByUserIDs は
	// 指定ユーザーたちのフォロワー・フォロー中の ID を userID ごとにまとめて返します。
	FindFollowerIDsByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]int, error)
	FindFollowingIDsByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]int, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}
