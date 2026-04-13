package repository

import (
	"context"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByName(ctx context.Context, name string) (*model.User, error)
}
