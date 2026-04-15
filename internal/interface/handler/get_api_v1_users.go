package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersHandler struct {
	InputPort usecase.GetUsersInputPort
}

func (h *GetApiV1UsersHandler) GetApiV1Users(ctx echo.Context, params openapi.GetApiV1UsersParams) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetUsersInput{
		Name: params.Name,
	})
	if err != nil {
		return err
	}

	response, err := openapi.NewGetUsersResponse(output.Users)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, response)
}
