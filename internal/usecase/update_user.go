package usecase

import (
	"context"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/domain/repository"
	"github.com/daisuke-harada/date-courses-go/internal/domain/service"
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

// Validate はユーザー更新の入力データをバリデーションします。
func (i *UpdateUserInput) Validate() error {
	var errs []string

	// gender: enum check
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

	// password: allow_nil（空なら検証スキップ）、指定時は6文字以上かつ確認一致
	if i.Password != "" {
		if len(i.Password) < 6 {
			errs = append(errs, "パスワードは6文字以上で入力してください")
		}
		if i.Password != i.PasswordConfirmation {
			errs = append(errs, "パスワード（確認）が一致しません")
		}
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	return nil
}

type UpdateUserOutput struct {
	UserWithRelations *model.UserWithRelations
}

type UpdateUserInteractor struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
}

func NewUpdateUserUsecase(
	userRepository repository.UserRepository,
	userService service.UserService,
) UpdateUserInputPort {
	return &UpdateUserInteractor{
		UserRepository: userRepository,
		UserService:    userService,
	}
}

func (i *UpdateUserInteractor) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// バリデーション
	if err := input.Validate(); err != nil {
		return nil, err
	}

	user, err := i.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, apperror.NotFound()
	}

	// パスワードが指定されている場合のみ更新（Rails の allow_nil: true に対応）
	if err := user.ApplyUpdate(input.Name, input.Email, input.Gender, input.Image, input.Password); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	if err := i.UserRepository.Update(ctx, user); err != nil {
		return nil, apperror.InternalServerError(err)
	}

	uwr, err := i.UserService.BuildUserWithRelations(ctx, user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{UserWithRelations: uwr}, nil
}

// NewUpdateUserInput builds UpdateUserInput from raw form string values.
// Invalid gender values are treated as empty and left for Validate() to report.
func NewUpdateUserInput(id int, name, email, genderStr, password, passwordConfirmation, imageStr string) (UpdateUserInput, error) {
	var gender model.Gender
	if genderStr != "" {
		g := model.Gender(genderStr)
		if g == model.GenderMale || g == model.GenderFemale {
			gender = g
		} else {
			gender = ""
		}
	}

	var image *string
	if imageStr != "" {
		image = &imageStr
	}

	return UpdateUserInput{
		ID:                   uint(id),
		Name:                 name,
		Email:                email,
		Gender:               gender,
		Image:                image,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	}, nil
}
