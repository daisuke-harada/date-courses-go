package handler

import (
	"net/http"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1UsersIdHandler struct {
	InputPort usecase.UpdateUserInputPort
}

func (h *PutApiV1UsersIdHandler) PutApiV1UsersId(ctx echo.Context, id int) error {
	var errs []string

	// multipart/form-data からフィールドを取得
	gender, err := NewModelGender(ctx.FormValue("gender"))
	if err != nil {
		errs = append(errs, "性別は「男性」または「女性」で入力してください")
	}

	input := usecase.UpdateUserInput{
		ID:                   uint(id),
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

	// password: allow_nil（空なら検証スキップ）、指定時は6文字以上かつ確認一致
	if input.Password != "" {
		if len(input.Password) < 6 {
			errs = append(errs, "パスワードは6文字以上で入力してください")
		}
		if input.Password != input.PasswordConfirmation {
			errs = append(errs, "パスワード（確認）が一致しません")
		}
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	uwr := output.UserWithRelations
	resp, err := openapi.BuildUserResponseBody(
		uwr.User,
		uwr.FollowerIDs,
		uwr.FollowingIDs,
		uwr.Courses,
		uwr.Reviews,
	)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
