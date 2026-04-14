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
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindFollowerIDsByUserID(ctx context.Context, userID uint) ([]int, error)
	FindFollowingIDsByUserID(ctx context.Context, userID uint) ([]int, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}
