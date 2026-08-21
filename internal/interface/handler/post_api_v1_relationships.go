package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1RelationshipsHandler struct {
	InputPort usecase.CreateRelationshipInputPort
}

func (h *PostApiV1RelationshipsHandler) PostApiV1Relationships(ctx echo.Context) error {
	// フォローする側はリクエストではなくトークンから決める。
	// current_user_id を信用すると他人になりすましてフォローさせられる。
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}

	input, err := usecase.NewCreateRelationshipInputFromStrings(currentUser.ID, ctx.FormValue("followed_user_id"))
	if err != nil {
		return err
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	resp, err := openapi.NewFollowResponseData(output)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusCreated, resp)
}
