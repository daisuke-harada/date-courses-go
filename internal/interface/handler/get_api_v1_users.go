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

	responses := make([]openapi.UserResponseBody, 0, len(output.Users))
	for _, uwr := range output.Users {
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
		responses = append(responses, resp)
	}

	return ctx.JSON(http.StatusOK, responses)
}
