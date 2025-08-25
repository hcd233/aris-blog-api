package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initAIRouter(r fiber.Router) {
	aiService := handler.NewAIHandler()
	aiRouter := r.Group("/ai", middleware.JwtMiddleware())
	{
		aiPromptRouter := aiRouter.Group("/prompt", middleware.LimitUserPermissionMiddleware("promptService", model.PermissionAdmin))
		{
			taskNameRouter := aiPromptRouter.Group("/:taskName", middleware.ValidateURIMiddleware(&protocol.TaskURI{}))
			{
				taskNameRouter.Get("/", middleware.ValidateParamMiddleware(&protocol.PaginateParam{}), aiService.HandleListPrompt)
				taskNameRouter.Post("/", middleware.ValidateBodyMiddleware(&protocol.CreatePromptBody{}), aiService.HandleCreatePrompt)
				taskNameRouter.Get("/latest", aiService.HandleGetLatestPrompt)
				promptVersionRouter := taskNameRouter.Group("/v:version", middleware.ValidateURIMiddleware(&protocol.PromptVersionURI{}))
				{
					promptVersionRouter.Get("/", aiService.HandleGetPrompt)
				}
			}
		}
		aiAppRouter := aiRouter.Group("/app")
		{
			creatorToolRouter := aiAppRouter.Group("/creator")
			{
				creatorToolRouter.Post(
					"/contentCompletion",
					middleware.ValidateBodyMiddleware(&protocol.GenerateContentCompletionBody{}),
					middleware.RedisLockMiddleware("contentCompletion", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateContentCompletion,
				)
				creatorToolRouter.Post(
					"/articleSummary",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleSummaryBody{}),
					middleware.RedisLockMiddleware("articleSummary", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateArticleSummary,
				)
				// creatorToolRouter.Post("/articleTranslation", aiService.GenerateArticleTranslationHandler)

			}
			readerToolRouter := aiAppRouter.Group("/reader")
			{
				readerToolRouter.Post(
					"/articleQA",
					middleware.ValidateBodyMiddleware(&protocol.GenerateArticleQABody{}),
					middleware.RedisLockMiddleware("articleQA", constant.CtxKeyUserID, 30*time.Second),
					aiService.HandleGenerateArticleQA,
				)
				// readerToolRouter.Post(
				// 	"/termExplaination",
				// 	middleware.ValidateBodyMiddleware(&protocol.GenerateTermExplainationBody{}),
				// 	aiService.HandleGenerateTermExplaination,
				// )
			}

		}
	}
}
