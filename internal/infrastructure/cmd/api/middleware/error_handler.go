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

	code, messages := resolveError(err)

	if writeErr := ctx.JSON(code, ErrorResponse{ErrorMessages: messages}); writeErr != nil {
		slog.Error("failed to write error response", "err", writeErr)
	}
}

// resolveError はエラーの種別を判定し、HTTPステータスコードとメッセージを返します。
func resolveError(err error) (code int, messages []string) {
	// ① apperror（usecase・handler どこからでも渡せる）
	if status, msgs, cause, ok := apperror.HTTPStatus(err); ok {
		if cause != nil {
			// cause（原因エラー）はサーバーログにのみ記録し、クライアントには返しません
			slog.Error("app error", "status", status, "messages", msgs, "cause", cause)
		} else {
			// 4xx 系など cause なし → Warn レベル
			slog.Warn("app error", "status", status, "messages", msgs)
		}
		return status, msgs
	}

	// ② echo.HTTPError（oapi-codegen のパラメータバインドエラーや Echo 内部から来る場合）
	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		msgs := echoMessagesToSlice(echoErr)
		slog.Error("echo http error", "status", echoErr.Code, "messages", msgs)
		return echoErr.Code, msgs
	}

	// ③ 予期しない error（バグや未ハンドルのケース）
	slog.Error("unexpected error", "err", err)
	return http.StatusInternalServerError, []string{http.StatusText(http.StatusInternalServerError)}
}

// echoMessagesToSlice は echo.HTTPError の Message フィールドを []string に変換します。
func echoMessagesToSlice(echoErr *echo.HTTPError) []string {
	switch m := echoErr.Message.(type) {
	case string:
		return []string{m}
	case []string:
		return m
	default:
		return []string{http.StatusText(echoErr.Code)}
	}
}
