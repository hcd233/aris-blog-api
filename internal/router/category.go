package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/handler"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

func initUserCategoryRouter(r *gin.RouterGroup) {
	categoryService := handler.NewCategoryService()

	r.GET("rootCategory", middleware.LimitUserPermissionMiddleware(model.PermissionCreator), categoryService.ListRootCategoriesHandler)
	categoryRouter := r.Group("/category", middleware.LimitUserPermissionMiddleware(model.PermissionCreator))
	{
		categoryRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), categoryService.CreateCategoryHandler)
	}

	categoryIDRouter := categoryRouter.Group(":categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
	{
		categoryIDRouter.GET("", categoryService.GetCategoryInfoHandler)
		categoryIDRouter.DELETE("", categoryService.DeleteCategoryHandler)
		categoryIDRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), categoryService.UpdateCategoryInfoHandler)
		categoryIDRouter.GET("subCategories", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryService.ListChildrenCategoriesHandler)
		categoryIDRouter.GET("subArticles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), categoryService.ListChildrenArticlesHandler)
	}
}
