// Package middleware 中间件
//
//	@update 2024-06-22 11:05:33
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/auth"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// JwtMiddleware JWT 中间件
//
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-16 05:35:57
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, protocol.Response{
				Code: protocol.CodeUnauthorized,
			})
			return
		}
		if isBearer := strings.HasPrefix(tokenString, "Bearer "); !isBearer {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code: protocol.CodeTokenVerifyError,
			})
			c.Abort()
			return
		}

		userID, err := auth.DecodeToken(tokenString[7:])
		if err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeTokenVerifyError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		user, err := model.QueryUserFieldsByID(userID, []string{"name", "permission"}, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUserNotFoundError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Set("userName", user.Name)
		c.Set("permission", user.Permission)
		c.Next()
	}
}
