package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1RelationshipsHandler struct {
	InputPort usecase.CreateRelationshipInputPort
}

func (h *PostApiV1RelationshipsHandler) PostApiV1Relationships(ctx echo.Context) error {
	input, err := usecase.NewCreateRelationshipInputFromStrings(ctx.FormValue("current_user_id"), ctx.FormValue("followed_user_id"))
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
