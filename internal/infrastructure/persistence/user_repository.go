package persistence

import (
	"context"
	"errors"
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

func (r *userRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("name = ?", name).Count(&count).Error; err != nil {
		slog.ErrorContext(ctx, "userRepository.ExistsByName failed", "err", err)
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		slog.ErrorContext(ctx, "userRepository.ExistsByEmail failed", "err", err)
		return false, err
	}
	return count > 0, nil
}
