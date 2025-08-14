package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"go.uber.org/zap"
)

// LogMiddleware 日志中间件
//
//	param logger *zap.Logger
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 21:21:46
func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		query := string(c.Request().URI().QueryString())

		err := c.Next()

		logger := logger.LoggerWithFiberContext(c)

		latency := time.Since(start)

		fields := []zap.Field{
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
			zap.String("latency", latency.String()),
			zap.String("req-content-type", c.Get("Content-Type")),
			zap.String("rsp-content-type", c.GetRespHeader("Content-Type")),
		}

		if err != nil {
			fields = append([]zap.Field{zap.Error(err)}, fields...)
			logger.Error("[FIBER] error", fields...)
		} else {
			logger.Info("[FIBER] info", fields...)
		}

		return err
	}
}
