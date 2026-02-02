package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type slogLogger struct {
	l   *slog.Logger
	ctx context.Context
}

func NewSlogLogger(handler slog.Handler) Logger {
	return &slogLogger{l: slog.New(handler), ctx: context.Background()}
}

func NewSlogLoggerFromEnv() Logger {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	// LOG_LEVEL は「どのレベル以上を出力するか」という閾値(フィルタ)です。
	// 例: LOG_LEVEL=info の場合 Debug() は捨てられ、Info/Warn/Error だけが出力されます。
	// ※Debug/Info/Warn/Error メソッド自体が「どのレベルで出すか」を決めます。
	level := parseLevel(os.Getenv("LOG_LEVEL"))

	opts := &slog.HandlerOptions{Level: level}

	switch format {
	case "", "text":
		return NewSlogLogger(slog.NewTextHandler(os.Stdout, opts))
	case "json":
		return NewSlogLogger(slog.NewJSONHandler(os.Stdout, opts))
	default:
		// 変な値でも落とさず text にフォールバック
		return NewSlogLogger(slog.NewTextHandler(os.Stdout, opts))
	}
}

func (s *slogLogger) With(ctx context.Context, args ...any) Logger {
	if ctx == nil {
		ctx = s.ctx
		if ctx == nil {
			ctx = context.Background()
		}
	}
	return &slogLogger{l: s.l.With(args...), ctx: ctx}
}

func (s *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	s.logCompat(ctx, slog.LevelDebug, msg, args...)
}
func (s *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	s.logCompat(ctx, slog.LevelInfo, msg, args...)
}
func (s *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	s.logCompat(ctx, slog.LevelWarn, msg, args...)
}
func (s *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	s.logCompat(ctx, slog.LevelError, msg, args...)
}

func (s *slogLogger) logCompat(ctx context.Context, level slog.Level, msg string, args ...any) {
	if ctx == nil {
		ctx = s.ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// 後方互換: Info("x=%v", v) のように、key/value ではなく printf 形式で渡されるケースがある。
	// これを検出して msg を整形し、slog の msg として流す。
	if looksLikePrintf(msg, args) {
		s.l.Log(ctx, level, fmt.Sprintf(msg, args...))
		return
	}
	s.l.Log(ctx, level, msg, args...)
}

func looksLikePrintf(msg string, args []any) bool {
	if len(args) == 0 {
		return false
	}
	// msg に % が含まれていれば printf 形式の可能性が高い
	return strings.Contains(msg, "%")
}

func parseLevel(s string) slog.Level {
	s = strings.ToLower(strings.TrimSpace(s))
	// レベルは下位 -> 上位 の順で定義（debug < info < warn < error）。
	// 環境変数の文字列に応じてハンドラ側の閾値(Level)を返す。
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
