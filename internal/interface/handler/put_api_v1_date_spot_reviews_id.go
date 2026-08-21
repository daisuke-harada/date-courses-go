package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/middleware"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotReviewsIdHandler struct {
	InputPort usecase.UpdateDateSpotReviewInputPort
}

func (h *PutApiV1DateSpotReviewsIdHandler) PutApiV1DateSpotReviewsId(ctx echo.Context, reviewID int) error {
	input, err := usecase.NewUpdateDateSpotReviewInput(
		reviewID,
		ctx.FormValue("date_spot_id"),
		ctx.FormValue("rate"),
		ctx.FormValue("content"),
	)
	if err != nil {
		return err
	}

	// 編集できるのは投稿者だけなので、操作主体はトークンから決める
	currentUser, err := middleware.RequireCurrentUser(ctx)
	if err != nil {
		return err
	}
	input.OperatorID = currentUser.ID

	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotReviewResponse(output.DateSpotReviews))
}
