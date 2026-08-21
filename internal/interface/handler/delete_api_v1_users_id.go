package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1UsersIdHandler struct {
	InputPort usecase.DeleteUserInputPort
}

func (h *DeleteApiV1UsersIdHandler) DeleteApiV1UsersId(ctx echo.Context, id int) error {
	// 退会できるのは本人だけなので、操作主体はトークンから決める
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}

	if err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteUserInput{
		ID:         uint(id),
		OperatorID: currentUser.ID,
	}); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
