package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initTokenRouter(r *gin.RouterGroup) {
	tokenService := service.NewTokenService()
	tokenRouter := r.Group("/token")
	{
		tokenRouter.POST(
			"refresh",
			middleware.RateLimiterMiddleware(config.JwtAccessTokenExpired/4, 2, "", protocol.CodeRefreshTokenRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.RefreshTokenBody{}),
			tokenService.RefreshTokenHandler,
		)
	}

}
