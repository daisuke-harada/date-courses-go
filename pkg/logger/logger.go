// This file does not contain a stray opening code fence.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	mu      sync.RWMutex
	appName string
	logger  *slog.Logger
	once    sync.Once
	logSink io.Writer
	// （現状はバッファリングを使わないため flush ハンドラは不要）
)

// request_id 用の context キー
type ctxKeyRequestID struct{}

// WithRequestID は context に request_id を付与する。
// 空文字列の場合は何もせずそのまま返す。
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKeyRequestID{}, requestID)
}

// RequestIDFromContext は context から request_id を取り出す。
// 存在しない、または空文字列の場合は false を返す。
func RequestIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v := ctx.Value(ctxKeyRequestID{})
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

// ログ設定:
//   - 出力形式は JSON のみサポートします（他フォーマットは不要なため除外）。
//   - LOG_LEVEL は newConfig() で決める（"debug" / "info" / "warn" / "error"）。
func Init(name string, debugFlag bool) {
	// sync.Once により初期化は1回のみ実行され、以降はリクエスト毎のロック不要
	once.Do(func() {
		l := newSlogLoggerFromConfig(name, debugFlag)
		logger = l
		// slog のパッケージレベルヘルパー (slog.Info 等) がこのロガーを使うように設定
		slog.SetDefault(logger)
	})
}

func Close() {
	mu.Lock()
	defer mu.Unlock()

	// デフォルトロガーを discard に切り替え、古いハンドラへの出力を防ぐ
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))

	logger = nil
	appName = ""
	debug = false
	logSink = nil
	once = sync.Once{}
}

func get() *slog.Logger {
	// 初期化が完了していることを保証する。once.Do はメモリ同期も提供するため、
	// 以降は logger の読み取りをロックなしで安全に行える。
	once.Do(func() {
		mu.RLock()
		n := appName
		d := debug
		mu.RUnlock()

		l := newSlogLoggerFromConfig(n, d)
		logger = l
		slog.SetDefault(logger)
	})
	return logger
}

// L は現在の *slog.Logger を返す。
// 初期化されていない場合は自動的に初期化する（遅延初期化）。
func L() *slog.Logger { return get() }

// With は現在のロガーに属性を追加した新しいロガーを返す。
func With(args ...any) *slog.Logger { return get().With(args...) }

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	get().Log(ctx, level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	get().LogAttrs(ctx, level, msg, attrs...)
}

func Enabled(ctx context.Context, level slog.Level) bool { return get().Enabled(ctx, level) }

func newSlogLoggerFromConfig(app string, debugFlag bool) *slog.Logger {
	setDebug(debugFlag)

	cfg := newConfig()

	// ログレベル: debugFlag が true なら強制的に debug
	level := parseLevel(cfg.Level)
	if debugFlag {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: level}

	// 出力先は常に stdout（バッファリングは無効）
	var sink io.Writer = os.Stdout

	// JSON ハンドラを常に使う
	var base slog.Handler = slog.NewJSONHandler(sink, opts)

	// sink を保存して Close() で使えるようにする
	mu.Lock()
	logSink = sink
	mu.Unlock()

	// app 属性をハンドラレベルで付与（全ログに自動的に含まれる）
	if app != "" {
		base = base.WithAttrs([]slog.Attr{slog.String("app", app)})
	}

	// context から request_id を自動で取り出してログに追加するハンドラでラップ
	h := &contextAttrsHandler{next: base}
	return slog.New(h)
}

// contextAttrsHandler は slog.Handler をラップし、context 内の属性（request_id 等）を
// ログレコードに自動追加する。
type contextAttrsHandler struct{ next slog.Handler }

func (h *contextAttrsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *contextAttrsHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := RequestIDFromContext(ctx); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}
	return h.next.Handle(ctx, r)
}

func (h *contextAttrsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextAttrsHandler{next: h.next.WithAttrs(attrs)}
}

func (h *contextAttrsHandler) WithGroup(name string) slog.Handler {
	return &contextAttrsHandler{next: h.next.WithGroup(name)}
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
