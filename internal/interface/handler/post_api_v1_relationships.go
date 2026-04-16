package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1RelationshipsHandler struct {
	InputPort usecase.CreateRelationshipInputPort
}

func (h *PostApiV1RelationshipsHandler) PostApiV1Relationships(ctx echo.Context) error {
	currentUserID, err := strconv.Atoi(ctx.FormValue("current_user_id"))
	if err != nil {
		return apperror.BadRequest("current_user_id は数値で指定してください")
	}

	followedUserID, err := strconv.Atoi(ctx.FormValue("followed_user_id"))
	if err != nil {
		return apperror.BadRequest("followed_user_id は数値で指定してください")
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.CreateRelationshipInput{
		CurrentUserID:  uint(currentUserID),
		FollowedUserID: uint(followedUserID),
	})
	if err != nil {
		return err
	}

	users := make([]openapi.UserResponseBody, 0, len(output.Users))
	for _, uwr := range output.Users {
		resp, err := openapi.BuildUserResponseBody(uwr.User, uwr.FollowerIDs, uwr.FollowingIDs, uwr.Courses, uwr.Reviews)
		if err != nil {
			return err
		}
		users = append(users, resp)
	}

	currentUser, err := openapi.BuildUserResponseBody(output.CurrentUser.User, output.CurrentUser.FollowerIDs, output.CurrentUser.FollowingIDs, output.CurrentUser.Courses, output.CurrentUser.Reviews)
	if err != nil {
		return err
	}

	followedUser, err := openapi.BuildUserResponseBody(output.FollowedUser.User, output.FollowedUser.FollowerIDs, output.FollowedUser.FollowingIDs, output.FollowedUser.Courses, output.FollowedUser.Reviews)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"users":         users,
		"current_user":  currentUser,
		"followed_user": followedUser,
	})
}
