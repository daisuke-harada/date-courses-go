package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type GetApiV1UsersHandler struct {}

func (h *GetApiV1UsersHandler) GetApiV1Users(ctx echo.Context) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}