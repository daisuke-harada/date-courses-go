package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler struct {
	InputPort usecase.DeleteRelationshipInputPort
}

func (h *DeleteApiV1RelationshipsCurrentUserIdOtherUserIdHandler) DeleteApiV1RelationshipsCurrentUserIdOtherUserId(ctx echo.Context, arg1 int, arg2 int) error {
	// 解除できるのは本人のフォローだけなので、操作主体はトークンから決める
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}

	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteRelationshipInput{
		UserID:     uint(arg1),
		FollowID:   uint(arg2),
		OperatorID: currentUser.ID,
	})
	if err != nil {
		return err
	}

	resp, err := openapi.NewUnFollowResponseData(output)
	if err != nil {
		return apperror.InternalServerError(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
