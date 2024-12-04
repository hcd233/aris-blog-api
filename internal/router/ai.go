package router

import (
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
				creatorToolRouter.POST("contentCompletion", middleware.ValidateBodyMiddleware(&protocol.GenerateContentCompletionBody{}), aiService.GenerateContentCompletionHandler)
				//creatorToolRouter.POST("articleSummary", aiService.GenerateArticleSummaryHandler)
				//creatorToolRouter.POST("articleTranslation", aiService.GenerateArticleTranslationHandler)

			}
			// readerToolRouter := aiAppRouter.Group("/reader")
			// {
			// 	readerToolRouter.POST("articleQA", aiService.GenerateArticleQAHandler)
			// 	readerToolRouter.POST("termExplaination", aiService.GenerateTermExplainationHandler)
			// }

		}
	}
}
