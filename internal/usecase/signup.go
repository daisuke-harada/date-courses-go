package usecase

import (
	"context"
	"regexp"
	"strings"

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

// emailRegex は Rails の validates_format_of :email と同等の正規表現です。
var emailRegex = regexp.MustCompile(`(?i)^[^@\s]+@(?:[-a-z0-9]+\.)+[a-z]{2,}$`)

func (i *SignupInteractor) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	// ─── バリデーション（Rails の validates に対応）─────────────────────
	var errs []string

	// name: presence, length(max:50)
	if strings.TrimSpace(input.Name) == "" {
		errs = append(errs, "名前を入力してください")
	} else if len(input.Name) > 50 {
		errs = append(errs, "名前は50文字以内で入力してください")
	}

	// email: presence, length(max:250), format
	if strings.TrimSpace(input.Email) == "" {
		errs = append(errs, "メールアドレスを入力してください")
	} else if len(input.Email) > 250 {
		errs = append(errs, "メールアドレスは250文字以内で入力してください")
	} else if !emailRegex.MatchString(input.Email) {
		errs = append(errs, "メールアドレスは正しい形式で入力してください")
	}

	// gender: presence
	if strings.TrimSpace(input.Gender) == "" {
		errs = append(errs, "性別を入力してください")
	}

	// password: presence, length(min:6)
	if input.Password == "" {
		errs = append(errs, "パスワードを入力してください")
	} else if len(input.Password) < 6 {
		errs = append(errs, "パスワードは6文字以上で入力してください")
	}

	// password_confirmation: match
	if input.Password != input.PasswordConfirmation {
		errs = append(errs, "パスワード（確認）が一致しません")
	}

	if len(errs) > 0 {
		return nil, apperror.UnprocessableEntity(errs...)
	}

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
