package handler

import (
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1LoginHandler struct {
	InputPort usecase.LoginInputPort
}

func (h *PostApiV1LoginHandler) PostApiV1Login(ctx echo.Context) error {
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.LoginInput{
		Name:     ctx.FormValue("name"),
		Password: ctx.FormValue("password"),
	})
	if err != nil {
		// apperror 型のエラーは CustomHTTPErrorHandler が適切なステータスコードで処理する
		return err
	}

	resp, err := openapi.NewLoginResponse(output.User)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, resp)
}
