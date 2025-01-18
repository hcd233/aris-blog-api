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
		userRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), userHandler.HandleQueryUser)
		userRouter.GET("me", userHandler.HandleGetCurUserInfo)
		userNameRouter := userRouter.Group("/:userName", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.GET("", userHandler.HandleGetUserInfo)
			userNameRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), userHandler.HandleUpdateInfo)

			initUserArticleRouter(userNameRouter)
			initUserCategoryRouter(userNameRouter)
			initUserOperationRouter(userNameRouter)
			initUserAssetRouter(userNameRouter)
		}

	}
}
