package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersHandler struct{}

func (h *GetApiV1UsersHandler) GetApiV1Users(ctx echo.Context, arg1 openapi.GetApiV1UsersParams) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
