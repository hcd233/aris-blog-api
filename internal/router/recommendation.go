package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// initRecommendationRouter 初始化推荐系统路由
//
//	author system
//	update 2025-01-19 12:00:00
func initRecommendationRouter(router fiber.Router) {
	// 获取数据库和Redis连接
	db := database.GetDBInstance(context.Background())
	redis := util.GetRedis()
	logger := util.GetLogger()

	// 创建推荐系统处理器
	recommendationHandler := handler.NewRecommendationHandler(db, redis, logger)

	// 推荐系统路由组
	recommendationGroup := router.Group("/recommendation")
	{
		// 用户行为上报（需要认证）
		recommendationGroup.Post("/behavior", middleware.Authenticate(), recommendationHandler.ReportBehavior)

		// 推荐接口（需要认证）
		recommendationGroup.Get("/articles", middleware.Authenticate(), recommendationHandler.RecommendArticles)
		recommendationGroup.Get("/tags", middleware.Authenticate(), recommendationHandler.RecommendTags)

		// 用户画像接口（需要认证）
		recommendationGroup.Get("/profile", middleware.Authenticate(), recommendationHandler.GetUserProfile)

		// 管理接口（需要管理员权限）
		adminGroup := recommendationGroup.Group("/admin", middleware.Authenticate(), middleware.RequirePermission("admin"))
		{
			adminGroup.Post("/train", recommendationHandler.TrainModel)
			adminGroup.Post("/profile/update", recommendationHandler.UpdateUserProfile)
		}
	}
}