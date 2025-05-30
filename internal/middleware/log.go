package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"go.uber.org/zap"
)

// LogMiddleware 日志中间件
//
//	param logger *zap.Logger
//	return gin.HandlerFunc
//	author centonhuang
//	update 2025-01-05 21:21:46
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger := logger.LoggerWithContext(c)

		latency := time.Since(start)

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("latency", latency.String()),
			zap.String("req-content-type", c.Request.Header.Get("Content-Type")),
			zap.String("rsp-content-type", c.Writer.Header().Get("Content-Type")),
		}

		if len(c.Errors) > 0 {
			fields = append([]zap.Field{zap.String("errors", c.Errors.String())}, fields...)
			logger.Error("[GIN] error", fields...)
		} else {
			logger.Info("[GIN] info", fields...)
		}
	}
}
