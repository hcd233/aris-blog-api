// Package middleware 中间件
//
//	@update 2024-06-22 11:05:33
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/auth"
	"github.com/hcd233/Aris-AI-go/internal/protocol"
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
				Code: protocol.CodeTokenVerifyError,
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
