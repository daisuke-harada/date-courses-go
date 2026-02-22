// This file does not contain a stray opening code fence.
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
	mu      sync.RWMutex
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

	v := ctx.Value(ctxKeyRequestID{})
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

// 環境変数:
//   - LOG_FORMAT: "text" (デフォルト) / "json"
//   - LOG_LEVEL: "debug" / "info" (デフォルト) / "warn" / "error"
//   - LOG_BUFFERED: "true" にすると 64KB バッファ付き出力を使う
func Init(name string, debugFlag bool) {
	// sync.Once により初期化は1回のみ実行され、以降はリクエスト毎のロック不要
	once.Do(func() {
		l := newSlogLoggerFromEnv(name, debugFlag)
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

	// アプリバッファがある場合はflushを呼ぶ
	// flush()を呼ぶとアプリバッファを下流に落とす、ファイルが出力先だったらバッファの情報もろともディスクに書き出される
	// ロック外で I/O 関連の終了処理を行う
	// 処理順と理由:
	// 1) logFlush (アプリ側の bufio.Flush 相当) を優先して呼ぶ
	//    - bufio.Writer 等のアプリ内バッファに溜まったデータを下流の Writer
	//      (例: os.Stdout や *os.File) に書き出す。
	//    - まずアプリバッファを吐き出すことで、ログがメモリに残るのを防ぐ。
	// 2) 下流が Sync() を実装している場合は Sync() を試す
	//    - これは *os.File の Sync() のように、カーネルのキャッシュから
	//      物理ディスクへ確実に書き込む（永続化）ための操作。重いので
	//      シャットダウン時にのみ行うのが普通。
	// 3) 次に汎用的な Flush() を試す（外部ライブラリ等の実装に対応）
	// 4) 最後に Close() を試してリソースを解放する（ソケットやファイル）
	//
	// いずれの I/O もブロッキングする可能性があるため、上ではロックを
	// 外してから実行している（ロック下では参照の取得と状態リセットのみ行う）。

	// すべてのクリーンアップを順に試す（早期 return をしない）
	// 1) アプリ側のバッファを落とす（logFlush）
	if flush != nil {
		if err := flush(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close app flush error: %v\n", err)
		}
	}

	if sink == nil {
		return
	}

	// 2) sink が Flush() error を持つ場合は呼ぶ
	if f, ok := sink.(interface{ Flush() error }); ok {
		if err := f.Flush(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close sink flush error: %v\n", err)
		}
	} else if f2, ok := sink.(interface{ Flush() }); ok {
		// Flush() を返さない型（例: 一部のストリーム）のためのフォールバック
		f2.Flush()
	}

	// 3) sink が Sync() を持つなら永続化を試みる
	if s, ok := sink.(interface{ Sync() error }); ok {
		if err := s.Sync(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close sync error: %v\n", err)
		}
	}

	// 4) 最後に Close() を試してリソースを解放
	if c, ok := sink.(interface{ Close() error }); ok {
		if err := c.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger.Close close error: %v\n", err)
		}
	}
}

func get() *slog.Logger {
	// 初期化が完了していることを保証する。once.Do はメモリ同期も提供するため、
	// 以降は logger の読み取りをロックなしで安全に行える。
	once.Do(func() {
		mu.RLock()
		n := appName
		d := debug
		mu.RUnlock()

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
