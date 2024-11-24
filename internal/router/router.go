package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/service"
)

// InitRouter initializes the router.
func InitRouter(r *gin.Engine) {
	pingService := service.NewPingService()

	r.GET("", pingService.PingHandler)

	v1Router := r.Group("/v1")
	{
		initTokenRouter(v1Router)
		initOauth2Router(v1Router)

		initUserRouter(v1Router)
		initTagRouter(v1Router)
		initArticleRouter(v1Router)

		initAIRouter(v1Router)
	}
}
