package router

import "github.com/gin-gonic/gin"

// StartupRouter is the router startup function.
func StartupRouter(r *gin.Engine) {
	group := r.Group("/")
	{
		group.GET("/", GetRootMessage)
	}
}
