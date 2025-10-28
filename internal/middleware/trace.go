package middleware

import (
    "context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/constant"
)

// TraceMiddleware 追踪中间件
//
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:30:00
func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-Id")

		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Locals(constant.CtxKeyTraceID, traceID)

		c.Set("X-Trace-Id", traceID)

        // 同步 TraceID 到 request context，供 Huma 使用
        uctx := c.UserContext()
        uctx = context.WithValue(uctx, constant.CtxKeyTraceID, traceID)
        c.SetUserContext(uctx)

		return c.Next()
	}
}
