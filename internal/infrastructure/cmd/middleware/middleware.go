package middleware

import (
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
	"github.com/labstack/echo/v4"
)

type MiddlewareInterface interface {
	LoggerMiddleware() echo.MiddlewareFunc
	// AuthMiddleware() echo.MiddlewareFunc
	// ValidationMiddleware() echo.MiddlewareFunc
}

type Middleware struct {
	log logger.LoggerInterface
}

func NewMiddleware(log logger.LoggerInterface) MiddlewareInterface {
	return &Middleware{log: log}
}
