package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
)

// LimitUserPermissionMiddleware 限制用户权限中间件
//
//	param serviceName string
//	param requiredPermission model.Permission
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:07:08
func LimitUserPermissionMiddleware(serviceName string, requiredPermission model.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permission := c.Locals(constant.CtxKeyPermission).(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			logger.WithFCtx(c).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			util.SendHTTPResponse(c, nil, protocol.ErrNoPermission)
			return c.Status(fiber.StatusForbidden).JSON(protocol.HTTPResponse{
				Error: protocol.ErrNoPermission.Error(),
			})
		}

		return c.Next()
	}
}

// PermissionMiddlewareForHuma 权限中间件 for Huma
//
//	param serviceName string
//	param requiredPermission model.Permission
//	return func(ctx huma.Context, next func(huma.Context))
//	author centonhuang
//	update 2025-01-05 21:00:00
func PermissionMiddlewareForHuma(serviceName string, requiredPermission model.Permission) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		permission := ctx.Value(constant.CtxKeyPermission).(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			logger.WithCtx(ctx.Context()).Info("[PermissionMiddlewareForHuma] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			ctx.SetStatus(fiber.StatusForbidden)
			return
		}

		next(ctx)
	}
}
