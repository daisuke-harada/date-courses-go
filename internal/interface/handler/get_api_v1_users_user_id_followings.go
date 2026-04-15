package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersUserIdFollowingsHandler struct {
	InputPort usecase.GetUserFollowingsInputPort
}

func (h *GetApiV1UsersUserIdFollowingsHandler) GetApiV1UsersUserIdFollowings(ctx echo.Context, userId int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetUserFollowingsInput{
		UserID: uint(userId),
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
