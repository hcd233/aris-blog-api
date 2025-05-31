package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initAIRouter(r *gin.RouterGroup) {
	aiService := handler.NewAIHandler()
	aiRouter := r.Group("/ai", middleware.JwtMiddleware())
	{
		aiPromptRouter := aiRouter.Group("/prompt", middleware.LimitUserPermissionMiddleware("promptService", model.PermissionAdmin))
		{
			taskNameRouter := aiPromptRouter.Group("/:taskName", middleware.ValidateURIMiddleware(&protocol.TaskURI{}))
			{
				taskNameRouter.GET("", middleware.ValidateParamMiddleware(&protocol.PageParam{}), aiService.HandleListPrompt)
				taskNameRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreatePromptBody{}), aiService.HandleCreatePrompt)
				taskNameRouter.GET("latest", aiService.HandleGetLatestPrompt)
				promptVersionRouter := taskNameRouter.Group("v:version", middleware.ValidateURIMiddleware(&protocol.PromptVersionURI{}))
				{
					promptVersionRouter.GET("", aiService.HandleGetPrompt)
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
					middleware.RedisLockMiddleware("contentCompletion", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateContentCompletion,
				)
				creatorToolRouter.POST(
					"articleSummary",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleSummaryBody{}),
					middleware.RedisLockMiddleware("articleSummary", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateArticleSummary,
				)
				// creatorToolRouter.POST("articleTranslation", aiService.GenerateArticleTranslationHandler)

			}
			readerToolRouter := aiAppRouter.Group("/reader")
			{
				readerToolRouter.POST(
					"articleQA",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleQABody{}),
					middleware.RedisLockMiddleware("articleQA", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateArticleQA,
				)
				// readerToolRouter.POST(
				// 	"termExplaination",
				// 	middleware.ValidateBodyMiddleware(&protocol.GenerateTermExplainationBody{}),
				// 	aiService.HandleGenerateTermExplaination,
				// )
			}

		}
	}
}
