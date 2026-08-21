package middleware

import (
	"sync"
	"time"

	"github.com/daisuke-harada/date-courses-go/internal/apperror"
	"github.com/labstack/echo/v4"
)

// rateLimitWindow は試行回数をカウントする時間窓です。
const rateLimitWindow = 1 * time.Minute

// rateLimitedPaths はレート制限をかけるパスです。
// 認証情報を試行できるエンドポイントに限定し、通常の閲覧は制限しません。
var rateLimitedPaths = map[string]struct{}{
	"/api/v1/login":  {},
	"/api/v1/signup": {},
}

// attemptCounter は IP ごとの試行回数を時間窓つきで数えます。
type attemptCounter struct {
	mu       sync.Mutex
	counts   map[string]int
	resetAt  time.Time
	window   time.Duration
	maxCount int
}

func newAttemptCounter(maxCount int, window time.Duration) *attemptCounter {
	return &attemptCounter{
		counts:   make(map[string]int),
		resetAt:  time.Now().Add(window),
		window:   window,
		maxCount: maxCount,
	}
}

// allow は key の試行を1回数え、上限内なら true を返します。
// 時間窓を過ぎたらカウンタ全体を作り直します。IP ごとにエントリを増やし続けないための措置です。
func (c *attemptCounter) allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().After(c.resetAt) {
		c.counts = make(map[string]int)
		c.resetAt = time.Now().Add(c.window)
	}

	c.counts[key]++
	return c.counts[key] <= c.maxCount
}

// LoginRateLimitMiddleware はログイン・新規登録への試行回数を IP ごとに制限します。
// 制限がないとパスワードのブルートフォースを何度でも試せるため、上限を設けます。
//
// カウンタはプロセス内のメモリに持ちます。Lambda では実行環境ごとに独立し、
// スケールアウトすると制限が緩くなる点に注意してください。
// 厳密な制限が必要になったら、API Gateway のスロットリングか外部ストアへ移します。
func LoginRateLimitMiddleware(maxAttemptsPerMinute int) echo.MiddlewareFunc {
	counter := newAttemptCounter(maxAttemptsPerMinute, rateLimitWindow)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if _, limited := rateLimitedPaths[ctx.Path()]; !limited {
				return next(ctx)
			}

			if !counter.allow(ctx.RealIP()) {
				return apperror.TooManyRequests()
			}

			return next(ctx)
		}
	}
}
