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
	// rate (optional): 型変換のみ。バリデーションは usecase.UpdateDateSpotReviewInput.Validate() が担う
	var rate *float64
	if rateStr := ctx.FormValue("rate"); rateStr != "" {
		r, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return apperror.BadRequest("rate は数値で指定してください")
		}
		rate = &r
	}

	// content (optional)
	var content *string
	if c := ctx.FormValue("content"); c != "" {
		content = &c
	}

	input := usecase.UpdateDateSpotReviewInput{
		ReviewID: uint(reviewID),
		Rate:     rate,
		Content:  content,
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotShowResponse(output.DateSpot, output.DateSpotReviews))
}
