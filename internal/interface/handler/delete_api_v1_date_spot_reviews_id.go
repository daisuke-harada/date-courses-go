package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type DeleteApiV1DateSpotReviewsIdHandler struct {
	InputPort usecase.DeleteDateSpotReviewInputPort
}

func (h *DeleteApiV1DateSpotReviewsIdHandler) DeleteApiV1DateSpotReviewsId(ctx echo.Context, arg1 int) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.DeleteDateSpotReviewInput{
		ReviewID: uint(arg1),
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, openapi.NewDateSpotShowResponse(output.DateSpot, output.DateSpotReviews))
}
