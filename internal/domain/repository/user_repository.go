package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
