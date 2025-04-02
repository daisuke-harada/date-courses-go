package middleware

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func (m *Middleware) Cors() echo.MiddlewareFunc {
	return echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	})
}
