package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initCategoryRouter(r *gin.RouterGroup) {
	categoryHandler := handler.NewCategoryHandler()

	categoryRouter := r.Group("/category",
		middleware.JwtMiddleware(),
		middleware.LimitUserPermissionMiddleware("categoryService", model.PermissionCreator))
	{
		categoryRouter.GET("root", categoryHandler.HandleGetRootCategories)
		categoryRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), categoryHandler.HandleCreateCategory)

		categoryIDRouter := categoryRouter.Group(":categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
		{
			categoryIDRouter.GET("", categoryHandler.HandleGetCategoryInfo)
			categoryIDRouter.DELETE("", categoryHandler.HandleDeleteCategory)
			categoryIDRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), categoryHandler.HandleUpdateCategoryInfo)
			categoryIDRouter.GET("subCategories", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenCategories)
			categoryIDRouter.GET("subArticles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryHandler.HandleListChildrenArticles)
		}
	}
}
