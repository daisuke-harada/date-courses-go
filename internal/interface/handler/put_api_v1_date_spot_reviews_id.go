package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PutApiV1DateSpotReviewsIdHandler struct {
	InputPort usecase.UpdateDateSpotReviewInputPort
}

func (h *PutApiV1DateSpotReviewsIdHandler) PutApiV1DateSpotReviewsId(ctx echo.Context, arg1 int) error {
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

	input := usecase.UpdateDateSpotReviewInput{
		ReviewID: uint(arg1),
		Rate:     rate,
		Content:  content,
	}
	_, err := h.InputPort.Execute(ctx.Request().Context(), input)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
