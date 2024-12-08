package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/service"
)

func initAIRouter(r *gin.RouterGroup) {
	aiService := service.NewAIService()
	aiRouter := r.Group("/ai", middleware.JwtMiddleware())
	{
		aiPromptRouter := aiRouter.Group("/prompt", middleware.LimitUserPermissionMiddleware(model.PermissionAdmin))
		{
			taskNameRouter := aiPromptRouter.Group("/:taskName", middleware.ValidateURIMiddleware(&protocol.TaskURI{}))
			{
				taskNameRouter.GET("", middleware.ValidateParamMiddleware(&protocol.PageParam{}), aiService.ListPromptHandler)
				taskNameRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreatePromptBody{}), aiService.CreatePromptHandler)
				taskNameRouter.GET("latest", aiService.GetLatestPromptHandler)
				promptVersionRouter := taskNameRouter.Group("v:version", middleware.ValidateURIMiddleware(&protocol.PromptVersionURI{}))
				{
					promptVersionRouter.GET("", aiService.GetPromptHandler)
				}
			}
		}
		aiAppRouter := aiRouter.Group("/app")
		{
			creatorToolRouter := aiAppRouter.Group("/creator")
			{
				creatorToolRouter.POST(
					"contentCompletion",
					middleware.ValidateBodyMiddleware(&protocol.GenerateContentCompletionBody{}),
					middleware.RedisLockMiddleware("contentCompletion", "userID", 30*time.Second),
					aiService.GenerateContentCompletionHandler,
				)
				creatorToolRouter.POST(
					"articleSummary",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleSummaryBody{}),
					middleware.RedisLockMiddleware("articleSummary", "userID", 30*time.Second),
					aiService.GenerateArticleSummaryHandler,
				)
				//creatorToolRouter.POST("articleTranslation", aiService.GenerateArticleTranslationHandler)

			}
			readerToolRouter := aiAppRouter.Group("/reader")
			{
				readerToolRouter.POST(
					"articleQA",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleQABody{}),
					middleware.RedisLockMiddleware("articleQA", "userID", 30*time.Second),
					aiService.GenerateArticleQAHandler,
				)
				// readerToolRouter.POST(
				// 	"termExplaination",
				// 	middleware.ValidateBodyMiddleware(&protocol.GenerateTermExplainationBody{}),
				// 	aiService.GenerateTermExplainationHandler,
				// )
			}

		}
	}
}
