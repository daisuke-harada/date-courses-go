package persistence

import (
	"context"
	"log/slog"

	model "github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.Create failed", "err", err)
		return err
	}
	slog.InfoContext(ctx, "userRepository.Create succeeded", "user_id", user.ID)
	return nil
}
