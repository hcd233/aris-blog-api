// Package middleware 中间件
//
//	update 2024-06-22 11:05:33
package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/auth"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// JwtMiddleware JWT 中间件
//
//	return gin.HandlerFunc
//	author centonhuang
//	update 2024-09-16 05:35:57
func JwtMiddleware() gin.HandlerFunc {
	db := database.GetDBInstance()
	dao := dao.GetUserDAO()
	jwtAccessTokenSvc := auth.GetJwtAccessTokenSigner()

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.LoggerWithContext(c).Error("[JwtMiddleware] token is empty")
			util.SendHTTPResponse(c, nil, protocol.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, err := jwtAccessTokenSvc.DecodeToken(tokenString)
		if err != nil {
			logger.LoggerWithContext(c).Error("[JwtMiddleware] failed to decode token", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrUnauthorized)
			c.Abort()
			return
		}

		user, err := dao.GetByID(db, userID, []string{"id", "name", "permission"}, []string{})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.LoggerWithContext(c).Error("[JwtMiddleware] user not found", zap.Uint("userID", userID))
				util.SendHTTPResponse(c, nil, protocol.ErrDataNotExists)
			} else {
				logger.LoggerWithContext(c).Error("[JwtMiddleware] failed to get user", zap.Uint("userID", userID), zap.Error(err))
				util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			}
			c.Abort()
			return
		}
		c.Set(constant.CtxKeyUserID, user.ID)
		c.Set(constant.CtxKeyUserName, user.Name)
		c.Set(constant.CtxKeyPermission, user.Permission)
		c.Next()
	}
}
