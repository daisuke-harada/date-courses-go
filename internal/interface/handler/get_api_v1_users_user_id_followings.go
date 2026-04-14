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
