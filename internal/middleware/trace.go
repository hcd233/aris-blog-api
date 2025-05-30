package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/constant"
)

// TraceMiddleware 追踪中间件
//
//	return gin.HandlerFunc
//	author centonhuang
//	update 2025-01-05 15:30:00
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-Id")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Set(constant.CtxKeyTraceID, traceID)

		c.Header("X-Trace-Id", traceID)

		c.Next()
	}
}
