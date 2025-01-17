// Package router 路由
package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/hcd233/aris-blog-api/docs"
	"github.com/hcd233/aris-blog-api/internal/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRouter 注册路由
//
//	param r *gin.Engine
//	author centonhuang
//	update 2025-01-04 15:32:40
func RegisterRouter(r *gin.Engine) {
	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pingService := handler.NewPingHandler()
	r.GET("", pingService.HandlePing)

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
