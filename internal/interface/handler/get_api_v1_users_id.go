package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersIdHandler struct {
	InputPort usecase.GetUserInputPort
}

func (h *GetApiV1UsersIdHandler) GetApiV1UsersId(ctx echo.Context, id int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.GetUserInput{
		ID: uint(id),
	})
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
