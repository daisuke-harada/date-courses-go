package logger

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
)

// Logger は、アプリケーション内で利用するロガーの抽象です。
//
// slog ベースの構造化ログを中心にしつつ、既存の fmt 形式(Info("x=%v", v)) も
// 互換のために残しています。
type Logger interface {
	With(ctx context.Context, args ...any) Logger
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

var (
	mu      sync.RWMutex
	current Logger
)

// Init はグローバルロガーを初期化します。
// 既に初期化済みなら何もしません。
func Init() error {
	// goroutineでの同時初期化を防止
	mu.Lock()
	defer mu.Unlock()

	// current が既に設定されていれば何もしない
	if current != nil {
		return nil
	}

	// LOGGER 環境変数に基づいてロガー実装を選択
	// 開発環境のためにswitchできるような形にしています
	impl := strings.ToLower(strings.TrimSpace(os.Getenv("LOGGER")))
	switch impl {
	case "slog":
		current = NewSlogLoggerFromEnv()
	default:
		return errors.New("unknown LOGGER value: " + impl)
	}

	return nil
}

// Close は後方互換のために残します。
// slog は基本的に明示的な close は不要ですが、将来の拡張に備えます。
func Close() {
	mu.Lock()
	defer mu.Unlock()
	current = nil
}

func get() Logger {
	mu.RLock()
	l := current
	mu.RUnlock()
	if l != nil {
		return l
	}
	mu.Lock()
	defer mu.Unlock()
	if current == nil {
		current = NewSlogLoggerFromEnv()
	}
	return current
}

// Get は現在のグローバルロガーを返します(未初期化ならデフォルトで初期化)。
func Get() Logger { return get() }

// Set はグローバルロガーを差し替えます(主にテスト用)。
func Set(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	current = l
}

// 後方互換 API
func NewLogger() error { return Init() }
func CloseLogger()     { Close() }

// 互換API: 既存の呼び出しは context.Background() を暗黙に使います。
// --- 推奨API: ctx を必ず渡す ---
func Info(ctx context.Context, msg string, args ...any)  { get().Info(ctx, msg, args...) }
func Error(ctx context.Context, msg string, args ...any) { get().Error(ctx, msg, args...) }
func Debug(ctx context.Context, msg string, args ...any) { get().Debug(ctx, msg, args...) }
func Warn(ctx context.Context, msg string, args ...any)  { get().Warn(ctx, msg, args...) }
