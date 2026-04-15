package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersUserIdFollowersHandler struct {
	InputPort usecase.GetUserFollowersInputPort
}

func (h *GetApiV1UsersUserIdFollowersHandler) GetApiV1UsersUserIdFollowers(ctx echo.Context, userId int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetUserFollowersInput{
		UserID: uint(userId),
	})
	if err != nil {
		return err
	}

	response, err := openapi.NewGetUsersResponse(output.Users)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}
