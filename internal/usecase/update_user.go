package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUserInputPort はユーザー更新ユースケースの入力ポートです。
type UpdateUserInputPort interface {
	Execute(context.Context, UpdateUserInput) (*UpdateUserOutput, error)
}

type UpdateUserInput struct {
	ID                   uint
	Name                 string
	Email                string
	Gender               model.Gender
	Image                *string
	Password             string
	PasswordConfirmation string
}

type UpdateUserOutput struct {
	UserWithRelations *UserWithRelations
}

type UpdateUserInteractor struct {
	UserRepository           repository.UserRepository
	CourseRepository         repository.CourseRepository
	DateSpotReviewRepository repository.DateSpotReviewRepository
}

func NewUpdateUserUsecase(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	dateSpotReviewRepository repository.DateSpotReviewRepository,
) UpdateUserInputPort {
	return &UpdateUserInteractor{
		UserRepository:           userRepository,
		CourseRepository:         courseRepository,
		DateSpotReviewRepository: dateSpotReviewRepository,
	}
}

func (i *UpdateUserInteractor) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Gender = input.Gender
	if input.Image != nil {
		user.Image = input.Image
	}

	// パスワードが指定されている場合のみ更新（Rails の allow_nil: true に対応）
	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
		user.PasswordDigest = string(hashed)
	}

	if err := i.UserRepository.Update(ctx, user); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	uwr, err := buildUserWithRelations(ctx, i.UserRepository, i.CourseRepository, i.DateSpotReviewRepository, user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{UserWithRelations: uwr}, nil
}
