package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// LimitUserPermissionMiddleware 限制用户权限中间件
// @param body interface{}
// @return gin.HandlerFunc
// @author centonhuang
// @update 2024-09-21 08:48:25
func LimitUserPermissionMiddleware(requiredPermission model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		permission := c.MustGet("permission").(model.Permission)
		if model.PermissionLevelMapping[permission] < model.PermissionLevelMapping[requiredPermission] {
			c.JSON(http.StatusForbidden, protocol.Response{
				Code:    protocol.CodeNotPermissionError,
				Message: fmt.Sprintf("Permission denied, required permission: %s, your permission: %s", requiredPermission, permission),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
