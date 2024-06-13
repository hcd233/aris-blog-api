package router

import "github.com/gin-gonic/gin"

// Router is the main router.
var Router = gin.Default()

// InitRouter initializes the router.
func InitRouter() {
	initRootRouter()
}

func initRootRouter() {
	rootGroup := Router.Group("/")
	{
		rootGroup.GET("/", rootHandler)
	}
}
