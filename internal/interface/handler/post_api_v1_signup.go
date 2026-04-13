package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1SignupHandler struct {
	InputPort usecase.SignupInputPort
}

func (h *PostApiV1SignupHandler) PostApiV1Signup(ctx echo.Context) error {
	// multipart/form-data からフィールドを取得
	input := usecase.SignupInput{
		Name:                 ctx.FormValue("name"),
		Email:                ctx.FormValue("email"),
		Gender:               ctx.FormValue("gender"),
		Password:             ctx.FormValue("password"),
		PasswordConfirmation: ctx.FormValue("password_confirmation"),
	}

	// image は任意フィールド
	if image := ctx.FormValue("image"); image != "" {
		input.Image = &image
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		// apperror 型のエラーは CustomHTTPErrorHandler が適切なステータスコードで処理する
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.NewSignupResponse(output.User))
}
