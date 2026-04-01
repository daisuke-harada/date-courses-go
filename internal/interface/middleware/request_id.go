package middleware

import (
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"github.com/labstack/echo/v4"
)

// RequestIDMiddleware injects the request id (from Echo's RequestID middleware)
// into the request's context so that slog.Context-based logging can pick it up.
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Try to get request id from header (set by middleware.RequestID)
		rid := c.Request().Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			// fallback to response header
			rid = c.Response().Header().Get(echo.HeaderXRequestID)
		}
		if rid != "" {
			// attach to request context
			r := c.Request()
			newCtx := logger.WithRequestID(r.Context(), rid)
			// create a shallow copy of request with the new context
			r = r.WithContext(newCtx)
			c.SetRequest(r)
		}
		return next(c)
	}
}
