package usecase

import (
	"context"
	"regexp"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
	jwtpkg "github.com/daisuke-harada/date-courses-go/internal/pkg/jwt"
)

// emailRegex は Rails の validates_format_of :email と同等の正規表現です。
var emailRegex = regexp.MustCompile(`(?i)^[^@\s]+@(?:[-a-z0-9]+\.)+[a-z]{2,}$`)

// SignupInputPort はサインアップユースケースの入力ポートです。
type SignupInputPort interface {
	Execute(context.Context, SignupInput) (*SignupOutput, error)
}

// SignupInput はサインアップの入力データです。
// OpenAPI の生成型に依存しないよう独自の型として定義します。
type SignupInput struct {
	Name                 string
	Email                string
	Gender               model.Gender
	Password             string
	PasswordConfirmation string
	Image                *string
}

// Validate はサインアップの入力データをバリデーションします。
func (i *SignupInput) Validate() error {
	var errs []string

	// gender: enum check (既に FormValue → model.Gender に変換済み)
	if i.Gender != model.GenderMale && i.Gender != model.GenderFemale {
		errs = append(errs, "性別は「男性」または「女性」で入力してください")
	}

	// name: presence, length(max:50)
	if strings.TrimSpace(i.Name) == "" {
		errs = append(errs, "名前を入力してください")
	} else if len(i.Name) > 50 {
		errs = append(errs, "名前は50文字以内で入力してください")
	}

	// email: presence, length(max:250), format
	if strings.TrimSpace(i.Email) == "" {
		errs = append(errs, "メールアドレスを入力してください")
	} else if len(i.Email) > 250 {
		errs = append(errs, "メールアドレスは250文字以内で入力してください")
	} else if !emailRegex.MatchString(i.Email) {
		errs = append(errs, "メールアドレスは正しい形式で入力してください")
	}

	// password: presence, length(min:6)
	if i.Password == "" {
		errs = append(errs, "パスワードを入力してください")
	} else if len(i.Password) < 6 {
		errs = append(errs, "パスワードは6文字以上で入力してください")
	}

	// password_confirmation: match
	if i.Password != i.PasswordConfirmation {
		errs = append(errs, "パスワード（確認）が一致しません")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	return nil
}

// SignupOutput はサインアップの出力データです。
type SignupOutput struct {
	User  *model.User
	Token string
}

type SignupInteractor struct {
	UserRepository repository.UserRepository
	AuthService    service.AuthService
	JWTSecretKey   string
}

func NewSignupUsecase(
	userRepository repository.UserRepository,
	authService service.AuthService,
	jwtSecret ...string,
) SignupInputPort {
	secret := ""
	if len(jwtSecret) > 0 {
		secret = jwtSecret[0]
	}
	return &SignupInteractor{
		UserRepository: userRepository,
		AuthService:    authService,
		JWTSecretKey:   secret,
	}
}

func (i *SignupInteractor) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	// バリデーション
	if err := input.Validate(); err != nil {
		return nil, err
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

	user := model.NewUser(input.Name, input.Email, input.Gender, input.Image, passwordDigest)

	if err := i.UserRepository.Create(ctx, user); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	out := &SignupOutput{User: user}
	// If a JWT secret is configured, issue a token so the user is logged in after signup
	if i.JWTSecretKey != "" {
		token, err := jwtpkg.Encode(user.ID, i.JWTSecretKey)
		if err != nil {
			return nil, apperror.InternalServerError(err)
		}
		out.Token = token
	}

	return out, nil
}
