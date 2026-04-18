package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1DateSpotReviewsHandler struct {
	InputPort usecase.CreateDateSpotReviewInputPort
}

func (h *PostApiV1DateSpotReviewsHandler) PostApiV1DateSpotReviews(ctx echo.Context) error {
	var errs []string

	// user_id (required)
	userIDStr := ctx.FormValue("user_id")
	if strings.TrimSpace(userIDStr) == "" {
		errs = append(errs, "ユーザーIDを入力してください")
	}
	userID, userIDErr := strconv.Atoi(userIDStr)
	if userIDErr != nil && strings.TrimSpace(userIDStr) != "" {
		errs = append(errs, "ユーザーIDは整数で入力してください")
	}

	// date_spot_id (required)
	dateSpotIDStr := ctx.FormValue("date_spot_id")
	if strings.TrimSpace(dateSpotIDStr) == "" {
		errs = append(errs, "デートスポットIDを入力してください")
	}
	dateSpotID, dateSpotIDErr := strconv.Atoi(dateSpotIDStr)
	if dateSpotIDErr != nil && strings.TrimSpace(dateSpotIDStr) != "" {
		errs = append(errs, "デートスポットIDは整数で入力してください")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	// rate (optional)
	var rate *float64
	if rateStr := ctx.FormValue("rate"); strings.TrimSpace(rateStr) != "" {
		r, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return apperror.UnprocessableEntity("評価は数値で入力してください")
		}
		rate = &r
	}

	// content (optional)
	var content *string
	if c := ctx.FormValue("content"); strings.TrimSpace(c) != "" {
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

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"review_id": int(output.ReviewID),
	})
}
