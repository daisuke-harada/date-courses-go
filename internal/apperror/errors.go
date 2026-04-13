package apperror

import (
	"errors"
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
// errors.As を使うことで、何重にWrapされていても appError を検出できます。
// error_handler.go で使用します。
func HTTPStatus(err error) (statusCode int, messages []string, cause error, ok bool) {
	var e *appError
	if !errors.As(err, &e) {
		return 0, nil, nil, false
	}
	return e.statusCode, e.messages, e.cause, true
}

// NotFound は 404 エラーを返します。
// msg を省略した場合は "リソースが見つかりません" がデフォルトメッセージになります。
func NotFound(msg ...string) error {
	message := "リソースが見つかりません"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &appError{
		statusCode: http.StatusNotFound,
		messages:   []string{message},
	}
}

// BadRequest は 400 エラーを返します。複数メッセージを渡せます（バリデーションエラー等）。
// msg を省略した場合は "リクエストが不正です" がデフォルトメッセージになります。
func BadRequest(msg ...string) error {
	messages := []string{"リクエストが不正です"}
	if len(msg) > 0 {
		messages = msg
	}
	return &appError{
		statusCode: http.StatusBadRequest,
		messages:   messages,
	}
}

// InternalServerError は 500 エラーを返します。
// cause に原因エラーを渡すと slog でログに記録されます（クライアントには返しません）。
// msg を省略した場合は "サーバーエラーが発生しました" がデフォルトメッセージになります。
func InternalServerError(cause error, msg ...string) error {
	message := "サーバーエラーが発生しました"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &appError{
		statusCode: http.StatusInternalServerError,
		messages:   []string{message},
		cause:      cause,
	}
}

// Unauthorized は 401 エラーを返します。
// msg を省略した場合は "認証が必要です" がデフォルトメッセージになります。
func Unauthorized(msg ...string) error {
	message := "認証が必要です"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &appError{
		statusCode: http.StatusUnauthorized,
		messages:   []string{message},
	}
}

// Forbidden は 403 エラーを返します。
// msg を省略した場合は "アクセスが禁止されています" がデフォルトメッセージになります。
func Forbidden(msg ...string) error {
	message := "アクセスが禁止されています"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &appError{
		statusCode: http.StatusForbidden,
		messages:   []string{message},
	}
}

// UnprocessableEntity は 422 エラーを返します。バリデーションエラー等で使用します。
// 複数メッセージを渡せます。
func UnprocessableEntity(msg ...string) error {
	messages := []string{"入力内容に誤りがあります"}
	if len(msg) > 0 {
		messages = msg
	}
	return &appError{
		statusCode: http.StatusUnprocessableEntity,
		messages:   messages,
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
