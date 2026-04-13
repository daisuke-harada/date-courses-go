package usecase

import (
	"context"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
)

// SignupInputPort はサインアップユースケースの入力ポートです。
type SignupInputPort interface {
	Execute(context.Context, SignupInput) (*SignupOutput, error)
}

// SignupInput はサインアップの入力データです。
// OpenAPI の生成型に依存しないよう独自の型として定義します。
type SignupInput struct {
	Name                 string
	Email                string
	Gender               string
	Password             string
	PasswordConfirmation string
	Image                *string
}

// SignupOutput はサインアップの出力データです。
type SignupOutput struct {
	User *model.User
}

type SignupInteractor struct {
	UserRepository repository.UserRepository
	AuthService    service.AuthService
}

func NewSignupUsecase(
	userRepository repository.UserRepository,
	authService service.AuthService,
) SignupInputPort {
	return &SignupInteractor{
		UserRepository: userRepository,
		AuthService:    authService,
	}
}

func (i *SignupInteractor) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	// ─── 一意性チェック（Rails の uniqueness に対応）─────────────────────
	nameExists, err := i.UserRepository.ExistsByName(ctx, input.Name)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if nameExists {
		return nil, apperror.UnprocessableEntity("名前はすでに存在します")
	}

	emailExists, err := i.UserRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}
	if emailExists {
		return nil, apperror.UnprocessableEntity("メールアドレスはすでに存在します")
	}

	// ─── パスワードハッシュ化（Rails の has_secure_password / bcrypt に対応）
	passwordDigest, err := i.AuthService.HashPassword(input.Password)
	if err != nil {
		return nil, apperror.InternalServerError(err)
	}

	// ─── ユーザー作成 ───────────────────────────────────────────────────
	user := &model.User{
		Name:           input.Name,
		Email:          input.Email,
		Gender:         input.Gender,
		Image:          input.Image,
		PasswordDigest: passwordDigest,
	}

	if err := i.UserRepository.Create(ctx, user); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	return &SignupOutput{User: user}, nil
}
