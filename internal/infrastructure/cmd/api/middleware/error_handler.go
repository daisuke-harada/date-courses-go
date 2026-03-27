package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/labstack/echo/v4"
)

// ErrorResponse はAPIのエラーレスポンスの形式を定義します。
// OpenAPI スキーマ (api/components/schemas/response/error.yaml) の ErrorResponse と対応しています。
type ErrorResponse struct {
	ErrorMessages []string `json:"errorMessages"`
}

// CustomHTTPErrorHandler は Echo のカスタムエラーハンドラーです。
// ハンドラーやミドルウェアから return された error をここで受け取り、
// 統一した ErrorResponse 形式で JSON レスポンスを返します。
//
// 呼び出しの流れ:
//
//	ハンドラ / usecase (return apperror.NotFound(...) など)
//	  → oapi-codegen ラッパー (api_server.gen.go) が return err
//	    → Echo ルーター → e.HTTPErrorHandler (= この関数)
func CustomHTTPErrorHandler(err error, ctx echo.Context) {
	// レスポンスが既に書き込み済みの場合は何もしない（二重レスポンス防止）
	if ctx.Response().Committed {
		return
	}

	var code = http.StatusInternalServerError
	var messages []string

	switch {
	// ① apperror（usecase・handler どこからでも渡せる）
	case func() bool { _, _, _, ok := apperror.HTTPStatus(err); return ok }():
		status, msgs, cause, _ := apperror.HTTPStatus(err)
		code = status
		messages = msgs
		if cause != nil {
			// cause（原因エラー）はサーバーログにのみ記録し、クライアントには返しません
			slog.Error("app error", "status", code, "messages", messages, "cause", cause)
		} else {
			// 4xx 系など cause なし → Warn レベル
			slog.Warn("app error", "status", code, "messages", messages)
		}

	// ② echo.HTTPError（oapi-codegen のパラメータバインドエラーや Echo 内部から来る場合）
	case errors.As(err, new(*echo.HTTPError)):
		var echoErr *echo.HTTPError
		errors.As(err, &echoErr)
		code = echoErr.Code
		switch m := echoErr.Message.(type) {
		case string:
			messages = []string{m}
		case []string:
			messages = m
		default:
			messages = []string{http.StatusText(code)}
		}
		slog.Error("echo http error", "status", code, "messages", messages)

	// ③ 予期しない error（バグや未ハンドルのケース）
	default:
		slog.Error("unexpected error", "err", err)
	}

	if writeErr := ctx.JSON(code, ErrorResponse{ErrorMessages: messages}); writeErr != nil {
		slog.Error("failed to write error response", "err", writeErr)
	}
}
