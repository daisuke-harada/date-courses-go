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

	resp, err := openapi.BuildCreateRelationshipResponse(output)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusCreated, resp)
}
