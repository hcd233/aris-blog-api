// Package middleware 中间件
//
//	update 2024-06-22 11:05:33
package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/jwt"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
)

// JwtMiddleware JWT 中间件
//
//	@return ctx huma.Context
//	@return next func(huma.Context)
//	@return func(ctx huma.Context, next func(huma.Context))
//	@author centonhuang
//	@update 2025-11-02 04:17:04
func JwtMiddleware() func(ctx huma.Context, next func(huma.Context)) {
	dao := dao.GetUserDAO()
	accessTokenSvc := jwt.GetAccessTokenSigner()

	return func(ctx huma.Context, next func(huma.Context)) {
		db := database.GetDBInstance(ctx.Context())

		tokenString := ctx.Header("Authorization")
		if tokenString == "" {
			ctx.SetStatus(fiber.StatusUnauthorized)
			return
		}
		userID, err := accessTokenSvc.DecodeToken(tokenString)
		if err != nil {
			ctx.SetStatus(fiber.StatusUnauthorized)
			return
		}
		user, err := dao.GetByID(db, userID, []string{"id", "name", "permission"}, []string{})
		if err != nil {
			ctx.SetStatus(fiber.StatusInternalServerError)
			return
		}
		ctx = huma.WithValue(ctx, constant.CtxKeyUserID, user.ID)
		ctx = huma.WithValue(ctx, constant.CtxKeyUserName, user.Name)
		ctx = huma.WithValue(ctx, constant.CtxKeyPermission, user.Permission)
		next(ctx)
	}
}
