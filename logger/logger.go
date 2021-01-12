package logger

import (
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const loggerKey = "logger"

func Middleware(logger *zap.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("X-Request-ID")
			l := logger.With((zap.String("x-request-id", id)))
			c.Set(loggerKey, l)
			return next(c)
		}
	}
}

// Extract extract zap logger from context
func Extract(c echo.Context) *zap.Logger {
	l, ok := c.Get(loggerKey).(*zap.Logger)
	if ok {
		return l
	}
	return zap.NewExample()
}
