package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

// Authenticate JWT认证中间件别名
//
//	author system
//	update 2025-01-19 12:00:00
func Authenticate() fiber.Handler {
	return JwtMiddleware()
}

// RequirePermission 权限检查中间件
//
//	param permission string
//	return fiber.Handler
//	author system
//	update 2025-01-19 12:00:00
func RequirePermission(permission string) fiber.Handler {
	var requiredPermission model.Permission
	switch permission {
	case "admin":
		requiredPermission = model.PermissionAdmin
	case "creator":
		requiredPermission = model.PermissionCreator
	case "reader":
		requiredPermission = model.PermissionReader
	default:
		requiredPermission = model.PermissionReader
	}
	
	return LimitUserPermissionMiddleware("recommendation", requiredPermission)
}