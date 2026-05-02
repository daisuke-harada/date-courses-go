package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1DateSpotReviewsHandler struct {
	InputPort usecase.CreateDateSpotReviewInputPort
}

func (h *PostApiV1DateSpotReviewsHandler) PostApiV1DateSpotReviews(ctx echo.Context) error {
	input, err := usecase.NewCreateDateSpotReviewInputFromStrings(
		ctx.FormValue("user_id"),
		ctx.FormValue("date_spot_id"),
		ctx.FormValue("rate"),
		ctx.FormValue("content"),
	)
	if err != nil {
		return err
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.NewDateSpotReviewResponse(output.DateSpotReviews))
}

// handler defers parsing to usecase.NewCreateDateSpotReviewInputFromStrings
