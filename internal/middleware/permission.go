package middleware

import (
	"github.com/gin-gonic/gin"
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
//	return gin.HandlerFunc
//	author centonhuang
//	update 2025-01-05 15:07:08
func LimitUserPermissionMiddleware(serviceName string, requiredPermission model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		permission := c.MustGet(constant.CtxKeyPermission).(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			logger.LoggerWithContext(c).Info("[LimitUserPermissionMiddleware] permission denied",
				zap.String("serviceName", serviceName),
				zap.String("requiredPermission", string(requiredPermission)),
				zap.String("permission", string(permission)))
			util.SendHTTPResponse(c, nil, protocol.ErrNoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}
