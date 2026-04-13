package handler

import (
	"net/http"
	"strings"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/daisuke-harada/date-courses-go/internal/interface/openapi"
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"github.com/labstack/echo/v4"
)

type PostApiV1LoginHandler struct {
	InputPort usecase.LoginInputPort
}

func (h *PostApiV1LoginHandler) PostApiV1Login(ctx echo.Context) error {
	name := ctx.FormValue("name")
	password := ctx.FormValue("password")

	// ─── バリデーション ────────────────────────────────────────────────
	var errs []string

	if strings.TrimSpace(name) == "" {
		errs = append(errs, "名前を入力してください")
	}
	if password == "" {
		errs = append(errs, "パスワードを入力してください")
	}

	if len(errs) > 0 {
		return apperror.UnprocessableEntity(errs...)
	}

	// ─── ユースケース実行 ──────────────────────────────────────────────
	output, err := h.InputPort.Execute(ctx.Request().Context(), usecase.LoginInput{
		Name:     name,
		Password: password,
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
