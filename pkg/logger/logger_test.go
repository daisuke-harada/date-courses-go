package logger_test

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

// --- test helpers (Step 1.2) ---

// envGuard は環境変数をテスト中だけ変更し、終了時に元に戻す。
func envGuard(t *testing.T, key, value string) {
	t.Helper()
	old, had := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("setenv %s: %v", key, err)
	}
	t.Cleanup(func() {
		if !had {
			_ = os.Unsetenv(key)
			return
		}
		_ = os.Setenv(key, old)
	})
}

// stdoutGuard は os.Stdout を一時的に差し替え、終了時に戻す。
func stdoutGuard(t *testing.T, w *os.File) {
	t.Helper()
	old := os.Stdout
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = old
	})
}

// stderrGuard は os.Stderr を一時的に差し替え、終了時に戻す。
// Close() が stderr へ出すメッセージをテスト出力に混ぜないために使う。
func stderrGuard(t *testing.T, w *os.File) {
	t.Helper()
	old := os.Stderr
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = old
	})
}

// recordingSink は Close()/Flush()/Sync() が呼ばれた順序や回数を記録できる writer。
// ※ Init が参照するのは os.Stdout（*os.File）なので、この型単体を os.Stdout に
// 直接差し込むことはできないが、Step 3 以降で「どのメソッドを呼びたいか」を
// 切り分けるために用意しておく。
type recordingSink struct {
	mu sync.Mutex

	writes  [][]byte
	calls   []string
	closed  bool
	flushed bool
	synced  bool

	flushErr error
	syncErr  error
	closeErr error
}

func (s *recordingSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, "Write")
	cp := make([]byte, len(p))
	copy(cp, p)
	s.writes = append(s.writes, cp)
	return len(p), nil
}

func (s *recordingSink) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, "Flush")
	s.flushed = true
	return s.flushErr
}

func (s *recordingSink) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, "Sync")
	s.synced = true
	return s.syncErr
}

func (s *recordingSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, "Close")
	s.closed = true
	return s.closeErr
}

func (s *recordingSink) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var all []byte
	for _, b := range s.writes {
		all = append(all, b...)
	}
	return string(all)
}

// flushNoErrSink は "Flush()"（error なし）を持つパターンの検証用。
type flushNoErrSink struct{ recordingSink }

func (s *flushNoErrSink) Flush() {
	s.recordingSink.mu.Lock()
	defer s.recordingSink.mu.Unlock()
	s.recordingSink.calls = append(s.recordingSink.calls, "Flush(noerr)")
	s.recordingSink.flushed = true
}

// errWriter は io.Writer としては成功するが、後片付け系メソッドでエラーを返す例。
// Step 5 で「エラーが出ても他のクリーンアップを続ける」ことの確認に使う。
type errWriter struct {
	io.Writer
}

var errSentinel = errors.New("sentinel")

// bufferFlushCloser はテスト専用の、Flush/Sync/Close を提供する writer。
// bufio.Writer の下流として使い、Close() が sink.Flush() (noerr) -> Sync -> Close
// を呼ぶことを検証する。
type bufferFlushCloser struct {
	buf        []byte
	mu         sync.Mutex
	flushCalls int
	syncCalls  int
	closeCalls int
}

func (b *bufferFlushCloser) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *bufferFlushCloser) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flushCalls++
}

func (b *bufferFlushCloser) Sync() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.syncCalls++
	return nil
}

func (b *bufferFlushCloser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeCalls++
	return nil
}

func (b *bufferFlushCloser) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return string(b.buf)
}

// TestWithRequestID は WithRequestID の基本動作を個別サブテストで検証する。
// 各サブテストで context の前後をログ出力して動作を見やすくしている。
func TestWithRequestID(t *testing.T) {
	t.Run("assigns request id", func(t *testing.T) {
		ctx := context.Background()
		rid := "req-123"

		beforeVal, beforeOk := logger.RequestIDFromContext(ctx)
		t.Logf("before: request_id ok=%v val=%q", beforeOk, beforeVal)

		got := logger.WithRequestID(ctx, rid)

		afterVal, afterOk := logger.RequestIDFromContext(got)
		t.Logf(" after: request_id ok=%v val=%q", afterOk, afterVal)

		if got == ctx {
			t.Fatalf("expected a new context when request id is set")
		}
		v, ok := logger.RequestIDFromContext(got)
		if !ok || v != rid {
			t.Fatalf("request_id mismatch: got ok=%v val=%q want ok=true val=%q", ok, v, rid)
		}
	})

	t.Run("empty does nothing", func(t *testing.T) {
		ctx := context.Background()
		rid := ""

		beforeVal, beforeOk := logger.RequestIDFromContext(ctx)
		t.Logf("before: request_id ok=%v val=%q", beforeOk, beforeVal)

		got := logger.WithRequestID(ctx, rid)

		afterVal, afterOk := logger.RequestIDFromContext(got)
		t.Logf(" after: request_id ok=%v val=%q", afterOk, afterVal)

		if got != ctx {
			t.Fatalf("expected the same context when request id is empty")
		}
		if afterOk {
			t.Fatalf("expected no request_id in context, but found %q", afterVal)
		}
	})

	t.Run("overwrites existing", func(t *testing.T) {
		ctx := logger.WithRequestID(context.Background(), "old")
		rid := "new"

		beforeVal, beforeOk := logger.RequestIDFromContext(ctx)
		t.Logf("before: request_id ok=%v val=%q", beforeOk, beforeVal)

		got := logger.WithRequestID(ctx, rid)

		afterVal, afterOk := logger.RequestIDFromContext(got)
		t.Logf(" after: request_id ok=%v val=%q", afterOk, afterVal)

		if got == ctx {
			t.Fatalf("expected a new context when overwriting request id")
		}
		v, ok := logger.RequestIDFromContext(got)
		if !ok || v != rid {
			t.Fatalf("request_id mismatch after overwrite: got ok=%v val=%q want ok=true val=%q", ok, v, rid)
		}
	})
}

// TestRequestIDFromContext は RequestIDFromContext の挙動を検証する。
// サブテストで nil / 空 / 値あり / WithRequestID で空文字の場合の挙動を確認する。
func TestRequestIDFromContext(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		var ctx context.Context = nil
		v, ok := logger.RequestIDFromContext(ctx)
		if ok || v != "" {
			t.Fatalf("expected (\"\", false) for nil context; got ok=%v val=%q", ok, v)
		}
	})

	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		v, ok := logger.RequestIDFromContext(ctx)
		if ok || v != "" {
			t.Fatalf("expected (\"\", false) for empty context; got ok=%v val=%q", ok, v)
		}
	})

	t.Run("present request id", func(t *testing.T) {
		ctx := logger.WithRequestID(context.Background(), "abc")
		v, ok := logger.RequestIDFromContext(ctx)
		if !ok || v != "abc" {
			t.Fatalf("expected (\"abc\", true); got ok=%v val=%q", ok, v)
		}
	})

	t.Run("empty request id via WithRequestID", func(t *testing.T) {
		base := context.Background()
		got := logger.WithRequestID(base, "")
		if got != base {
			t.Fatalf("expected same context when WithRequestID called with empty string")
		}
		v, ok := logger.RequestIDFromContext(got)
		if ok || v != "" {
			t.Fatalf("expected (\"\", false) after empty WithRequestID; got ok=%v val=%q", ok, v)
		}
	})
}

func TestLogger_InitClose_UnbufferedWritesToStdout(t *testing.T) {
	// LOG_BUFFERED=false で Init->Log->Close を行い、stdout に書き込まれることを確認する。
	// ※ここでは Close の Sync/Close の呼び出し順までは見ず、「出力される」ことを重視する。

	// 環境変数はテスト内で明示して固定する
	envGuard(t, "LOG_BUFFERED", "false")
	envGuard(t, "LOG_FORMAT", "text")
	envGuard(t, "LOG_LEVEL", "debug")

	// stdout をパイプに差し替えて出力を回収する
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	stdoutGuard(t, w)
	t.Cleanup(func() {
		_ = w.Close()
		_ = r.Close()
	})

	// Close() が pipe に対して Sync() を試みると stderr にエラーが出るため、吸い込む
	derrR, derrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(stderr): %v", err)
	}
	stderrGuard(t, derrW)
	t.Cleanup(func() {
		_ = derrW.Close()
		_ = derrR.Close()
	})

	// Init してログを出す
	logger.Close() // テスト間の影響を避けるため、念のためリセット
	logger.Init("test-app", false)
	logger.Log(context.Background(), slog.LevelInfo, "hello", "k", "v")
	logger.Close()

	// 書き込み側を閉じて読み取りを完了させる
	_ = w.Close()
	got, readErr := io.ReadAll(r)
	if readErr != nil {
		t.Fatalf("read stdout pipe: %v", readErr)
	}
	if len(got) == 0 {
		t.Fatalf("expected stdout output, got empty")
	}
	// text handler の出力に msg が含まれることを最低限確認
	if !strings.Contains(string(got), "hello") {
		t.Fatalf("expected log message to contain %q, got %q", "hello", string(got))
	}
}

func TestLogger_InitClose_BufferedFlushesToStdout(t *testing.T) {
	// LOG_BUFFERED=true で Init->Log->Close を行い、bufio の flush によって
	// stdout（下流）へデータが流れることを確認する。
	// ※Close 内部で logFlush を呼ぶこと（= bufio.Writer.Flush 相当）が目的。

	// 環境変数はテスト内で明示して固定する
	envGuard(t, "LOG_BUFFERED", "true")
	envGuard(t, "LOG_FORMAT", "text")
	envGuard(t, "LOG_LEVEL", "debug")

	// stdout をパイプに差し替えて出力を回収する
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	stdoutGuard(t, w)
	t.Cleanup(func() {
		_ = w.Close()
		_ = r.Close()
	})

	// Init してログを出す
	logger.Close() // テスト間の影響を避けるため、念のためリセット
	logger.Init("test-app", false)

	// bufio.Writer を噛ませている場合、(ログ出力) 時点では下流へ落ちない可能性があるため
	// Close() を呼んで flush させた上で読み取る。
	logger.Log(context.Background(), slog.LevelInfo, "hello-buffered", "k", "v")
	logger.Close()

	// 書き込み側を閉じて読み取りを完了させる
	_ = w.Close()
	got, readErr := io.ReadAll(r)
	if readErr != nil {
		t.Fatalf("read stdout pipe: %v", readErr)
	}
	if len(got) == 0 {
		t.Fatalf("expected stdout output, got empty")
	}
	if !strings.Contains(string(got), "hello-buffered") {
		t.Fatalf("expected log message to contain %q, got %q", "hello-buffered", string(got))
	}
}

func TestLogger_Close_CallsSinkFlushNoErrSyncClose(t *testing.T) {
	// Close() が sink の interface 判定を使って
	// - Flush() (error なし)
	// - Sync() error
	// - Close() error
	// を順に試すことを検証する。
	//
	// ここでは Init() を使わず、テスト側から logSink/logFlush を仕込む。
	// （exported API ではないが、Close のロジック確認のための白箱テスト）

	// テスト間の影響を避ける
	logger.Close()

	// bufio.NewWriter が吐き出す先を自前で用意し、logFlush は bufio 側を使用。
	under := &bufferFlushCloser{}
	bw := bufio.NewWriterSize(under, 16)

	// パッケージ変数 logSink/logFlush を直接セットできないため、
	// Init を経由してセットする必要がある。
	// ただし Init は sink を os.Stdout に固定するため、ここでは Close の
	// sink 側ロジックのみを検証するために、Close() の前処理を模擬する。
	//
	// 具体的には、logger パッケージは Close 内で logFlush/logSink をスナップショットするので、
	// その値を作るには Init/newSlogLoggerFromEnv が必要。
	// → 外部テスト (package logger_test) では内部変数を触れないため、
	// 本テストは「Flush(noerr)/Sync/Close を持つ sink を Close に通す」ことを
	// 直接は検証できない。
	//
	// 代替案: os.Stdout を *os.File 以外に差し込めない制約があるため、
	// Step 5 では「Flush() error のみ（bufio.Writer）」の検証に留め、
	// noerr Flush の検証は sink 注入ポイントを設計変更した上で行う。
	_ = bw
	_ = under
	// このテストは設計上の理由で現在は実装できないため、明示的にスキップする。
	t.Skip("Close の sink 注入ポイントがないため、外部パッケージテストから Flush() (noerr) の呼び出しを直接検証できない")
}
