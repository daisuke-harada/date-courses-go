package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type PostApiV1UserNameSearchHandler struct {}

func (h *PostApiV1UserNameSearchHandler) PostApiV1UserNameSearch(ctx echo.Context) error {
	// TODO: Implement your logic here
	// Example: return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}