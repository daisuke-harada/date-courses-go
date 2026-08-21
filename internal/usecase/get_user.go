package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// GetUserInputPort はユーザー単体取得ユースケースの入力ポートです。
type GetUserInputPort interface {
	Execute(context.Context, GetUserInput) (*GetUserOutput, error)
}

type GetUserInput struct {
	ID uint
	// ViewerID は閲覧しているユーザーの ID です。未ログインの場合は 0 になります。
	ViewerID uint
}

type GetUserOutput struct {
	UserWithRelations *model.UserWithRelations
}

type GetUserInteractor struct {
	UserRepository   repository.UserRepository
	CourseRepository repository.CourseRepository
	UserService      service.UserService
}

func NewGetUserUsecase(
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	userService service.UserService,
) GetUserInputPort {
	return &GetUserInteractor{
		UserRepository:   userRepository,
		CourseRepository: courseRepository,
		UserService:      userService,
	}
}

func (i *GetUserInteractor) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	uwr, err := i.UserService.BuildUserWithRelations(ctx, user)
	if err != nil {
		return nil, err
	}

	// 非公開コースを見られるのは、本人が自分のマイページを開いたときだけ。
	// それ以外の経路では UserService が公開コースしか取得しない。
	if input.ViewerID != 0 && input.ViewerID == input.ID {
		courses, err := i.CourseRepository.FindAllByUserID(ctx, input.ID)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
		uwr.Courses = courses
	}

	return &GetUserOutput{UserWithRelations: uwr}, nil
}
