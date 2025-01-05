package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
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
			initUserTagRouter(userNameRouter)
			initUserOperationRouter(userNameRouter)
			initUserAssetRouter(userNameRouter)
		}

	}
}
