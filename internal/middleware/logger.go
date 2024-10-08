package middleware

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func formatDuration(d time.Duration) string {
	if d.Seconds() >= 1 {
		return fmt.Sprintf("%.2f s", d.Seconds())
	} else if d.Milliseconds() >= 1 {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else if d.Microseconds() >= 1 {
		return fmt.Sprintf("%d µs", d.Microseconds())
	}
	return fmt.Sprintf("%d ns", d.Nanoseconds())
}

// LoggerMiddleware 日志中间件
//
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-16 01:05:49
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := lo.Must1(url.QueryUnescape(c.Request.URL.RawQuery))
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
				zap.String("latency", formatDuration(latency)),
			)
		}
	}
}
