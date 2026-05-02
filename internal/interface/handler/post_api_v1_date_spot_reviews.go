package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1DateSpotReviewsHandler struct {
	InputPort usecase.CreateDateSpotReviewInputPort
}

func (h *PostApiV1DateSpotReviewsHandler) PostApiV1DateSpotReviews(ctx echo.Context) error {
	// user_id: 型変換のみ。バリデーションは usecase.CreateDateSpotReviewInput.Validate() が担う
	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
	if err != nil {
		return apperror.BadRequest("user_id は整数で指定してください")
	}

	// date_spot_id: 型変換のみ
	dateSpotID, err := strconv.Atoi(ctx.FormValue("date_spot_id"))
	if err != nil {
		return apperror.BadRequest("date_spot_id は整数で指定してください")
	}

	// rate (optional)
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

	input := usecase.CreateDateSpotReviewInput{
		UserID:     uint(userID),
		DateSpotID: uint(dateSpotID),
		Rate:       rate,
		Content:    content,
	}
	output, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, openapi.NewDateSpotShowResponse(output.DateSpot, output.DateSpotReviews))
}
