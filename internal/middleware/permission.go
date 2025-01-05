package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/util"
	"go.uber.org/zap"
)

// LimitUserPermissionMiddleware 限制用户权限中间件
//
//	@param serviceName string
//	@param requiredPermission model.Permission
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2025-01-05 15:07:08
func LimitUserPermissionMiddleware(serviceName string, requiredPermission model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		permission := c.MustGet("permission").(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			logger.Logger.Info("[LimitUserPermissionMiddleware] permission denied",
				zap.Uint("userID", userID),
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
