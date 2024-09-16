package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/logger"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
//
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-16 01:05:49
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Logger.Error(e.Error())
			}
		} else {
			logger.Logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.Duration("latency", latency),
			)
		}
	}
}
