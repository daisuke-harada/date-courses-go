package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1UsersIdHandler struct {
	InputPort usecase.UpdateUserInputPort
}

func (h *PutApiV1UsersIdHandler) PutApiV1UsersId(ctx echo.Context, id int) error {
	// multipart/form-data からフィールドを取得
	gender, err := NewModelGender(ctx.FormValue("gender"))
	if err != nil {
		gender = ""
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

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	resp, err := openapi.NewUserWithRelationsResponse(output.UserWithRelations)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
