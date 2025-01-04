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
	tokenRouter := r.Group("/token")
	{
		tokenRouter.POST(
			"refresh",
			middleware.RateLimiterMiddleware(config.JwtAccessTokenExpired/4, 2, "refreshToken", "userID", protocol.CodeRefreshTokenRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.RefreshTokenBody{}),
			tokenHandler.HandleRefreshToken,
		)
	}
}

