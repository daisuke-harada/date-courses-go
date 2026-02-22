package logger

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	mu      sync.Mutex
	appName string
	debug   bool
	logger  *slog.Logger
	once    sync.Once
	// バッファ付き出力を使う場合に保持する sink と flush 関数
	// LOG_BUFFERED=true のとき bufio.Writer を使い、Close() で Flush する
	logSink  io.Writer
	logFlush func() error
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

	// context から request_id を取り出すためにキーで検索する。
	// - キーにはパッケージローカルな型 ctxKeyRequestID の空値を使う（型を分けることで他パッケージとの衝突を防ぐ）。
	// - ctx.Value はキーが存在しなければ nil を返すので、下で型アサーションと空文字列チェックを行う。
	v := ctx.Value(ctxKeyRequestID{})
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

// Init はアプリケーション名とデバッグモードを設定し、ロガーを初期化する。
// 複数回呼ばれても初期化は1回だけ実行される（sync.Once）。
// 初期化後は slog.SetDefault が呼ばれるため、パッケージレベルの
// slog.Info/Error 等がこのロガーを使うようになる。
//
// 環境変数:
//   - LOG_FORMAT: "text" (デフォルト) / "json"
//   - LOG_LEVEL: "debug" / "info" (デフォルト) / "warn" / "error"
//   - LOG_BUFFERED: "true" にすると 64KB バッファ付き出力を使う
func Init(name string, debugFlag bool) {
	mu.Lock()
	appName = name
	debug = debugFlag
	mu.Unlock()

	// sync.Once により初期化は1回のみ実行され、以降はリクエスト毎のロック不要
	once.Do(func() {
		mu.Lock()
		n := appName
		d := debug
		mu.Unlock()

		l := newSlogLoggerFromEnv(n, d)
		logger = l
		// slog のパッケージレベルヘルパー (slog.Info 等) がこのロガーを使うように設定
		slog.SetDefault(logger)
	})
}

// Close はロガーの状態をリセットし、バッファ付き出力の場合は Flush を実行する。
// 主にテスト用途で再初期化を可能にするために用意している。
//
// 注意: os.Exit を使う場合は defer が実行されないため、Close は呼ばれない。
// LOG_BUFFERED=true のときは明示的に Close を呼ぶか、defer logger.Close() を推奨。
func Close() {
	// ロック下で flush/sink を取得し、状態をリセット
	mu.Lock()
	flush := logFlush
	sink := logSink
	logger = nil
	appName = ""
	debug = false
	logSink = nil
	logFlush = nil
	// once をリセットして再初期化可能にする（主にテスト用）
	once = sync.Once{}
	mu.Unlock()

	// ロック外で flush を実行（ブロックを避けるため）
	if flush != nil {
		if err := flush(); err != nil {
			// logger はすでに nil なので stderr に直接出力
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close flush error: %v\n", err)
		}
		return
	}

	// flush がない場合、sink が Sync/Flush/Close を実装していれば呼ぶ
	type syncer interface{ Sync() error }
	if s, ok := sink.(syncer); ok {
		if err := s.Sync(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close sync error: %v\n", err)
		}
		return
	}
	type flusher interface{ Flush() error }
	if f, ok := sink.(flusher); ok {
		if err := f.Flush(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close flush error: %v\n", err)
		}
		return
	}
	type closer interface{ Close() error }
	if c, ok := sink.(closer); ok {
		if err := c.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close close error: %v\n", err)
		}
		return
	}
}

func get() *slog.Logger {
	// 初期化が完了していることを保証する。once.Do はメモリ同期も提供するため、
	// 以降は logger の読み取りをロックなしで安全に行える。
	once.Do(func() {
		mu.Lock()
		n := appName
		d := debug
		mu.Unlock()

		l := newSlogLoggerFromEnv(n, d)
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
func newSlogLoggerFromEnv(app string, debugFlag bool) *slog.Logger {
	// 環境変数 LOG_FORMAT から出力形式を取得（text / json）
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))

	// 環境変数 LOG_LEVEL からログレベルを取得、debugFlag が true なら強制的に debug
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	if debugFlag {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: level}

	// 出力先の選択: LOG_BUFFERED=true ならバッファ付き、それ以外は直接 stdout
	var sink io.Writer = os.Stdout
	var flush func() error
	if strings.ToLower(strings.TrimSpace(os.Getenv("LOG_BUFFERED"))) == "true" {
		bw := bufio.NewWriterSize(os.Stdout, 64*1024)
		sink = bw
		flush = bw.Flush
	}

	var base slog.Handler
	switch format {
	case "", "text":
		base = slog.NewTextHandler(sink, opts)
	case "json":
		base = slog.NewJSONHandler(sink, opts)
	default:
		// 不正な値の場合も text にフォールバック（落とさない）
		base = slog.NewTextHandler(sink, opts)
	}

	// sink/flush を保存して Close() で使えるようにする
	mu.Lock()
	logSink = sink
	logFlush = flush
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
