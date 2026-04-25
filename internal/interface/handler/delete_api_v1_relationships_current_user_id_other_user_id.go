package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler struct {
	InputPort usecase.DeleteRelationshipInputPort
}

func (h *DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler) DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx echo.Context, arg1 int, arg2 int) error {
	err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteRelationshipInput{
		UserID:   uint(arg1),
		FollowID: uint(arg2),
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
