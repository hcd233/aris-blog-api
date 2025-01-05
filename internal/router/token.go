package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

func initTokenRouter(r *gin.RouterGroup) {
	tokenHandler := handler.NewTokenHandler()

	tokenRouter := r.Group("/token", middleware.JwtMiddleware())
	{
		tokenRouter.POST(
			"refresh",
			middleware.RateLimiterMiddleware("refreshToken", "userID", config.JwtAccessTokenExpired/4, 2),
			middleware.ValidateBodyMiddleware(&protocol.RefreshTokenBody{}),
			tokenHandler.HandleRefreshToken,
		)
	}
}
