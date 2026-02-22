package logger_test

import (
	"context"
	"testing"

	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

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
