package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
	"gorm.io/gorm"
)

type LoginInputPort interface {
	Execute(context.Context, LoginInput) (*LoginOutput, error)
}

type LoginInput struct {
	Name     string
	Password string
}

func (i *LoginInput) Validate() error {
	var errs []string

	if strings.TrimSpace(i.Name) == "" {
		errs = append(errs, "名前を入力してください")
	}
	if i.Password == "" {
		errs = append(errs, "パスワードを入力してください")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	return nil
}

type LoginOutput struct {
	User  *model.User
	Token string
}

type JWTSecretKey string

type LoginInteractor struct {
	UserRepository repository.UserRepository
	AuthService    service.AuthService
	JWTSecretKey   JWTSecretKey
}

func NewLoginUsecase(
	userRepository repository.UserRepository,
	authService service.AuthService,
	jwtSecretKey JWTSecretKey,
) LoginInputPort {
	return &LoginInteractor{
		UserRepository: userRepository,
		AuthService:    authService,
		JWTSecretKey:   jwtSecretKey,
	}
}

func (i *LoginInteractor) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := i.UserRepository.FindByName(ctx, input.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized(
				"認証に失敗しました。",
				"正しい名前・パスワードを入力し直すか、新規登録を行ってください。",
			)
		}
		return nil, apperror.InternalServerError(err)
	}

	if !i.AuthService.CheckPassword(user.PasswordDigest, input.Password) {
		return nil, apperror.Unauthorized(
			"認証に失敗しました。",
			"正しい名前・パスワードを入力し直すか、新規登録を行ってください。",
		)
	}

	token, err := jwtpkg.Encode(user.ID, string(i.JWTSecretKey))
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &LoginOutput{User: user, Token: token}, nil
}
