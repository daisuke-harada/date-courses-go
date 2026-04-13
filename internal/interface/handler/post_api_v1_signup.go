package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

// emailRegex は Rails の validates_format_of :email と同等の正規表現です。
var emailRegex = regexp.MustCompile(`(?i)^[^@\s]+@(?:[-a-z0-9]+\.)+[a-z]{2,}$`)

type PostApiV1SignupHandler struct {
	InputPort usecase.SignupInputPort
}

func (h *PostApiV1SignupHandler) PostApiV1Signup(ctx echo.Context) error {
	// ─── バリデーション（Rails の validates に対応）─────────────────────
	var errs []string

	// multipart/form-data からフィールドを取得
	gender, err := NewModelGender(ctx.FormValue("gender"))
	if err != nil {
		errs = append(errs, "性別は「男性」または「女性」で入力してください")
	}

	input := usecase.SignupInput{
		Name:                 ctx.FormValue("name"),
		Email:                ctx.FormValue("email"),
		Gender:               gender,
		Password:             ctx.FormValue("password"),
		PasswordConfirmation: ctx.FormValue("password_confirmation"),
	}

	// image は任意フィールド
	if image := ctx.FormValue("image"); image != "" {
		input.Image = &image
	}

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
		return apperror.UnprocessableEntity(errs...)
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		// apperror 型のエラーは CustomHTTPErrorHandler が適切なステータスコードで処理する
		return err
	}

	response, err := openapi.NewSignupResponse(output.User)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, response)
}

// toModelGender は FormValue の文字列を model.Gender に変換します。
// 「男性」「女性」以外の値は error を返します。
func NewModelGender(s string) (model.Gender, error) {
	switch model.Gender(s) {
	case model.GenderMale:
		return model.GenderMale, nil
	case model.GenderFemale:
		return model.GenderFemale, nil
	default:
		return "", fmt.Errorf("invalid gender value: %q", s)
	}
}
