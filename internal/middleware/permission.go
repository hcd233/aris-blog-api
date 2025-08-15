package middleware

import (
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
			logger.LoggerWithFiberContext(c).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			util.SendHTTPResponse(c, nil, protocol.ErrNoPermission)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		return c.Next()
	}
}
