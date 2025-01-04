package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

func initUserRouter(r *gin.RouterGroup) {
	userService := handler.NewUserService()

	userRouter := r.Group("/user", middleware.JwtMiddleware())
	{
		userRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), userService.QueryUserHandler)
		userRouter.GET("me", userService.GetMyInfoHandler)
		userNameRouter := userRouter.Group("/:userName", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.GET("", userService.GetUserInfoHandler)
			userNameRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), userService.UpdateInfoHandler)

			initUserArticleRouter(userNameRouter)
			initUserCategoryRouter(userNameRouter)
			initUserTagRouter(userNameRouter)
			initUserOperationRouter(userNameRouter)
			initUserAssetRouter(userNameRouter)
		}

	}
}
