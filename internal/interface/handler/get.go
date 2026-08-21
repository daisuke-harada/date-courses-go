package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetHandler struct{}

func (h *GetHandler) Get(ctx echo.Context) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
