package handler

import (
	"fmt"
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1SignupHandler struct {
	InputPort usecase.SignupInputPort
}

func (h *PostApiV1SignupHandler) PostApiV1Signup(ctx echo.Context) error {
	// multipart/form-data からフィールドを取得

	input, err := usecase.NewSignupInput(
		ctx.FormValue("name"),
		ctx.FormValue("email"),
		ctx.FormValue("gender"),
		ctx.FormValue("password"),
		ctx.FormValue("password_confirmation"),
		ctx.FormValue("image"),
	)
	if err != nil {
		return err
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

	response, err := openapi.NewSignupResponse(output.User, output.Token)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, response)
}

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
