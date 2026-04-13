package usecase

import (
	"context"
	"errors"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	"gorm.io/gorm"
)

// LoginInputPort はログインユースケースの入力ポートです。
type LoginInputPort interface {
	Execute(context.Context, LoginInput) (*LoginOutput, error)
}

// LoginInput はログインの入力データです。
type LoginInput struct {
	Name     string
	Password string
}

// LoginOutput はログインの出力データです。
type LoginOutput struct {
	User *model.User
}

type LoginInteractor struct {
	UserRepository repository.UserRepository
	AuthService    service.AuthService
}

func NewLoginUsecase(
	userRepository repository.UserRepository,
	authService service.AuthService,
) LoginInputPort {
	return &LoginInteractor{
		UserRepository: userRepository,
		AuthService:    authService,
	}
}

func (i *LoginInteractor) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// name でユーザーを検索
	user, err := i.UserRepository.FindByName(ctx, input.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Rails と同じメッセージで 401 を返す
			return nil, apperror.Unauthorized(
				"認証に失敗しました。",
				"正しい名前・パスワードを入力し直すか、新規登録を行ってください。",
			)
		}
		return nil, apperror.InternalServerError(err)
	}

	// bcrypt でパスワードを検証（Rails の user.authenticate(password) 相当）
	if !i.AuthService.CheckPassword(user.PasswordDigest, input.Password) {
		return nil, apperror.Unauthorized(
			"認証に失敗しました。",
			"正しい名前・パスワードを入力し直すか、新規登録を行ってください。",
		)
	}

	return &LoginOutput{User: user}, nil
}
