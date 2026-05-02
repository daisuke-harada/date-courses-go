package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotReviewsIdHandler struct {
	InputPort usecase.UpdateDateSpotReviewInputPort
}

func (h *PutApiV1DateSpotReviewsIdHandler) PutApiV1DateSpotReviewsId(ctx echo.Context, reviewID int) error {
	dateSpotID, err := strconv.Atoi(ctx.FormValue("date_spot_id"))
	if err != nil {
		return apperror.BadRequest("date_spot_id は整数で指定してください")
	}

	var rate *float64
	if rateStr := ctx.FormValue("rate"); rateStr != "" {
		r, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return apperror.BadRequest("rate は数値で指定してください")
		}
		rate = &r
	}

	var content *string
	if c := ctx.FormValue("content"); c != "" {
		content = &c
	}

	input := usecase.UpdateDateSpotReviewInput{
		ReviewID:   uint(reviewID),
		DateSpotID: uint(dateSpotID),
		Rate:       rate,
		Content:    content,
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotReviewResponse(output.DateSpotReviews))
}
