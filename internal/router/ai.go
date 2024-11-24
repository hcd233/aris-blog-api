package router

import "github.com/gin-gonic/gin"

func initAIRouter(r *gin.RouterGroup) {
	// 	aiRouter := r.Group("/ai", middleware.JwtMiddleware())
	// 	{
	// 		initAICreatorToolRouter(aiRouter)
	// 		initAIReaderToolRouter(aiRouter)
	// 	}
	// }

	// func initAICreatorToolRouter(r *gin.RouterGroup) {
	// 	creatorToolRouter := r.Group("/creator")
	// 	{
	// 		creatorToolRouter.POST("summary", ai.GetSummaryHandler)
	// 		creatorToolRouter.POST("translation", ai.GetTranslationHandler)

	// 	}
	// }

	//	func initAIReaderToolRouter(r *gin.RouterGroup) {
	//		readerToolRouter := r.Group("/reader")
	//		{
	//			readerToolRouter.POST("articleQA", ai.GetArticleQAHandler)
	//			readerToolRouter.POST("creatorIntro", ai.GetCreatorIntroHandler)
	//		}
}
