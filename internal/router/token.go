package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initTokenRouter(r *gin.RouterGroup) {
	tokenHandler := handler.NewTokenHandler()

	tokenRouter := r.Group("/token")
	{
		tokenRouter.POST(
			"refresh",
			middleware.RateLimiterMiddleware("refreshToken", "", config.JwtAccessTokenExpired/4, 2),
			middleware.ValidateBodyMiddleware(&protocol.RefreshTokenBody{}),
			tokenHandler.HandleRefreshToken,
		)
	}
}
