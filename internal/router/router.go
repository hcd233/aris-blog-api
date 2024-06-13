package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/router/v1/oauth2"
)

// Router is the main router.
var Router = gin.Default()

// InitRouter initializes the router.
func InitRouter(r *gin.Engine) {
	initRootRouter(r)
	initV1Router(r)
}

func initRootRouter(r *gin.Engine) {
	rootGroup := r.Group("/")
	{
		rootGroup.GET("/", handleRoot)
	}
}

func initV1Router(r *gin.Engine) {
	v1Router := r.Group("/v1")
	{
		oauth2.InitOauth2Router(v1Router)
	}
}
