// Package middleware 中间件
//
//	update 2024-06-22 11:05:33
package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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
//	return fiber.Handler
//	author centonhuang
//	update 2024-09-16 05:35:57
func JwtMiddleware() fiber.Handler {
	dao := dao.GetUserDAO()
	jwtAccessTokenSvc := auth.GetJwtAccessTokenSigner()

	return func(c *fiber.Ctx) error {
		db := database.GetDBInstanceFromFiber(c)

		tokenString := c.Get("Authorization")
		if tokenString == "" {
			logger.WithFCtx(c).Error("[JwtMiddleware] token is empty")
			util.SendHTTPResponse(c, nil, protocol.ErrUnauthorized)
			return c.Status(fiber.StatusUnauthorized).JSON(protocol.HTTPResponse[any]{
				Error: protocol.ErrUnauthorized.Error(),
			})
		}

		userID, err := jwtAccessTokenSvc.DecodeToken(tokenString)
		if err != nil {
			logger.WithFCtx(c).Error("[JwtMiddleware] failed to decode token", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrUnauthorized)
			return c.Status(fiber.StatusUnauthorized).JSON(protocol.HTTPResponse[any]{
				Error: protocol.ErrUnauthorized.Error(),
			})
		}

		user, err := dao.GetByID(db, userID, []string{"id", "name", "permission"}, []string{})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.WithFCtx(c).Error("[JwtMiddleware] user not found", zap.Uint("userID", userID))
				util.SendHTTPResponse(c, nil, protocol.ErrDataNotExists)
			} else {
				logger.WithFCtx(c).Error("[JwtMiddleware] failed to get user", zap.Uint("userID", userID), zap.Error(err))
				util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse[any]{
				Error: protocol.ErrInternalError.Error(),
			})
		}
		c.Locals(constant.CtxKeyUserID, user.ID)
		c.Locals(constant.CtxKeyUserName, user.Name)
		c.Locals(constant.CtxKeyPermission, user.Permission)
		return c.Next()
	}
}
