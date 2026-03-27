package apperror

import (
	"fmt"
	"net/http"
)

// appError はアプリケーション独自のエラー型です（外部には非公開）。
// HTTPステータスコード・クライアントへ返すメッセージ・原因エラー（Wrap）を保持します。
type appError struct {
	statusCode int
	messages   []string
	cause      error // Wrap元のerror。slogでのログ出力に使用し、クライアントには返しません。
}

// Error は error インターフェースの実装です。
func (e *appError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("status=%d messages=%v: %v", e.statusCode, e.messages, e.cause)
	}
	return fmt.Sprintf("status=%d messages=%v", e.statusCode, e.messages)
}

// Unwrap により errors.Is / errors.As で原因エラーを辿れるようにします。
func (e *appError) Unwrap() error {
	return e.cause
}

// HTTPStatus は error が appError かどうか判定し、HTTPステータス・メッセージ・原因エラーを返します。
// error_handler.go で使用します。
func HTTPStatus(err error) (statusCode int, messages []string, cause error, ok bool) {
	e, ok := err.(*appError)
	if !ok {
		return 0, nil, nil, false
	}
	return e.statusCode, e.messages, e.cause, true
}

// --- コンストラクタ（返り値は error 型） ---
// usecase・handler どこからでも普通の error として扱えます。

// NotFound は 404 エラーを返します。
func NotFound(msg string) error {
	return &appError{
		statusCode: http.StatusNotFound,
		messages:   []string{msg},
	}
}

// BadRequest は 400 エラーを返します。複数メッセージを渡せます（バリデーションエラー等）。
func BadRequest(messages ...string) error {
	return &appError{
		statusCode: http.StatusBadRequest,
		messages:   messages,
	}
}

// InternalServerError は 500 エラーを返します。
// cause に原因エラーを渡すと slog でログに記録されます（クライアントには返しません）。
func InternalServerError(cause error, msg string) error {
	return &appError{
		statusCode: http.StatusInternalServerError,
		messages:   []string{msg},
		cause:      cause,
	}
}

// Unauthorized は 401 エラーを返します。
func Unauthorized(msg string) error {
	return &appError{
		statusCode: http.StatusUnauthorized,
		messages:   []string{msg},
	}
}

// Forbidden は 403 エラーを返します。
func Forbidden(msg string) error {
	return &appError{
		statusCode: http.StatusForbidden,
		messages:   []string{msg},
	}
}

// Wrap は任意のステータスコードで既存の error をラップします。
func Wrap(cause error, statusCode int, msg string) error {
	return &appError{
		statusCode: statusCode,
		messages:   []string{msg},
		cause:      cause,
	}
}
