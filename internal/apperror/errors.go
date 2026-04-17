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

// newAppError は appError を生成する内部ヘルパーです。
// msg が指定されなければ defaultMsg を使います。cause は slog 用に保持されます。
func newAppError(statusCode int, defaultMsg string, cause error, msg []string) *appError {
	message := defaultMsg
	if len(msg) > 0 {
		message = msg[0]
	}
	return &appError{
		statusCode: statusCode,
		messages:   []string{message},
		cause:      cause,
	}
}

// newAppErrors は複数メッセージを扱う内部ヘルパーです。
func newAppErrors(statusCode int, defaultMsgs []string, cause error, msgs []string) *appError {
	messages := defaultMsgs
	if len(msgs) > 0 {
		messages = msgs
	}
	return &appError{
		statusCode: statusCode,
		messages:   messages,
		cause:      cause,
	}
}

// ---------- 404 Not Found ----------

// NotFound は 404 エラーを返します。
// msg を省略した場合は "リソースが見つかりません" がデフォルトメッセージになります。
func NotFound(msg ...string) error {
	return newAppError(http.StatusNotFound, "リソースが見つかりません", nil, msg)
}

// NotFoundWithCause は cause を slog 用にラップした 404 エラーを返します。
// msg を省略した場合はデフォルトメッセージが使われます。
func NotFoundWithCause(cause error, msg ...string) error {
	return newAppError(http.StatusNotFound, "リソースが見つかりません", cause, msg)
}

// ---------- 400 Bad Request ----------

// BadRequest は 400 エラーを返します。
// msg を省略した場合は "リクエストが不正です" がデフォルトメッセージになります。
func BadRequest(msg ...string) error {
	return newAppErrors(http.StatusBadRequest, []string{"リクエストが不正です"}, nil, msg)
}

// BadRequestWithCause は cause を slog 用にラップした 400 エラーを返します。
func BadRequestWithCause(cause error, msg ...string) error {
	return newAppErrors(http.StatusBadRequest, []string{"リクエストが不正です"}, cause, msg)
}

// ---------- 401 Unauthorized ----------

// Unauthorized は 401 エラーを返します。
// msg を省略した場合は "認証が必要です" がデフォルトメッセージになります。
func Unauthorized(msg ...string) error {
	return newAppErrors(http.StatusUnauthorized, []string{"認証が必要です"}, nil, msg)
}

// UnauthorizedWithCause は cause を slog 用にラップした 401 エラーを返します。
func UnauthorizedWithCause(cause error, msg ...string) error {
	return newAppErrors(http.StatusUnauthorized, []string{"認証が必要です"}, cause, msg)
}

// ---------- 403 Forbidden ----------

// Forbidden は 403 エラーを返します。
// msg を省略した場合は "アクセスが禁止されています" がデフォルトメッセージになります。
func Forbidden(msg ...string) error {
	return newAppError(http.StatusForbidden, "アクセスが禁止されています", nil, msg)
}

// ForbiddenWithCause は cause を slog 用にラップした 403 エラーを返します。
func ForbiddenWithCause(cause error, msg ...string) error {
	return newAppError(http.StatusForbidden, "アクセスが禁止されています", cause, msg)
}

// ---------- 422 Unprocessable Entity ----------

// UnprocessableEntity は 422 エラーを返します。
// msg を省略した場合は "入力内容に誤りがあります" がデフォルトメッセージになります。
func UnprocessableEntity(msg ...string) error {
	return newAppErrors(http.StatusUnprocessableEntity, []string{"入力内容に誤りがあります"}, nil, msg)
}

// UnprocessableEntityWithCause は cause を slog 用にラップした 422 エラーを返します。
func UnprocessableEntityWithCause(cause error, msg ...string) error {
	return newAppErrors(http.StatusUnprocessableEntity, []string{"入力内容に誤りがあります"}, cause, msg)
}

// ---------- 500 Internal Server Error ----------

// InternalServerError は 500 エラーを返します。
// cause に原因エラーを渡すと slog でログに記録されます（クライアントには返しません）。
// msg を省略した場合は "サーバーエラーが発生しました" がデフォルトメッセージになります。
func InternalServerError(cause error, msg ...string) error {
	return newAppError(http.StatusInternalServerError, "サーバーエラーが発生しました", cause, msg)
}

// ---------- 汎用ラッパー ----------

// Wrap は任意のステータスコードで既存の error をラップします。
func Wrap(cause error, statusCode int, msg string) error {
	return &appError{
		statusCode: statusCode,
		messages:   []string{msg},
		cause:      cause,
	}
}
