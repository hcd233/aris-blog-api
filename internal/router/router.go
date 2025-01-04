// Package router 路由
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
)

// RegisterRouter 注册路由
//	@param r *gin.Engine 
//	@author centonhuang 
//	@update 2025-01-04 15:32:40 
func RegisterRouter(r *gin.Engine) {
	pingService := handler.NewPingService()

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
