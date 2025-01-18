package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

func initUserRouter(r *gin.RouterGroup) {
	userHandler := handler.NewUserHandler()

	userRouter := r.Group("/user", middleware.JwtMiddleware())
	{
		userRouter.GET("current", userHandler.HandleGetCurUserInfo)
		userRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), userHandler.HandleUpdateInfo)
		userNameRouter := userRouter.Group("/:userID", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.GET("", userHandler.HandleGetUserInfo)
		}

	}
}
