package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/middleware"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/router/v1/article"
	article_version "github.com/hcd233/Aris-blog/internal/router/v1/article_version"
	"github.com/hcd233/Aris-blog/internal/router/v1/category"
	"github.com/hcd233/Aris-blog/internal/router/v1/comment"
	"github.com/hcd233/Aris-blog/internal/router/v1/like"
	"github.com/hcd233/Aris-blog/internal/router/v1/oauth2"
	"github.com/hcd233/Aris-blog/internal/router/v1/tag"
	"github.com/hcd233/Aris-blog/internal/router/v1/user"
)

// InitRouter initializes the router.
func InitRouter(r *gin.Engine) {
	rootGroup := r.Group("/")
	{
		rootGroup.GET("/", RootHandler)
	}

	v1Router := r.Group("/v1")
	{
		initOauth2Router(v1Router)
		initTagRouter(v1Router)
		initArticleRouter(v1Router)
		initUserRouter(v1Router)
	}
}

func initOauth2Router(r *gin.RouterGroup) {
	oauth2Group := r.Group("/oauth2")
	{
		githubRouter := oauth2Group.Group("/github")
		{
			githubRouter.GET("/login", oauth2.GithubLoginHandler)
			githubRouter.GET("/callback", oauth2.GithubCallbackHandler)
		}
	}
}

func initTagRouter(r *gin.RouterGroup) {
	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tag.ListTagsHandler)
	tagRouter := r.Group("/tag", middleware.JwtMiddleware())
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tag.QueryTagHandler)
		tagRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateTagBody{}), tag.CreateTagHandler)
		tagSlugRouter := tagRouter.Group("/:tagSlug", middleware.ValidateURIMiddleware(&protocol.TagURI{}))
		{
			tagSlugRouter.GET("", tag.GetTagInfoHandler)
			tagSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateTagBody{}), tag.UpdateTagHandler)
			tagSlugRouter.DELETE("", tag.DeleteTagHandler)
		}
	}
}

func initArticleRouter(r *gin.RouterGroup) {
	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), article.ListArticlesHandler)
	articleRouter := r.Group("/article", middleware.JwtMiddleware())
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), article.QueryArticleHandler)
	}
}

func initUserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("/user", middleware.JwtMiddleware())
	{
		userRouter.GET("/", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), user.QueryUserHandler)

		userNameRouter := userRouter.Group("/:userName", middleware.ValidateURIMiddleware(&protocol.UserURI{}))
		{
			userNameRouter.GET("", user.GetUserInfoHandler)
			userNameRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateUserBody{}), user.UpdateInfoHandler)

			initUserArticleRouter(userNameRouter)
			initUserCategoryRouter(userNameRouter)
			initUserTagRouter(userNameRouter)
			initUserOperationRouter(userNameRouter)
		}

	}
}

func initUserArticleRouter(r *gin.RouterGroup) {
	r.GET("articles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), article.ListUserArticlesHandler)
	articleRouter := r.Group("/article")
	{
		articleRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), article.QueryUserArticleHandler)
		articleRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateArticleBody{}), article.CreateArticleHandler)
	}

	articleSlugRouter := articleRouter.Group("/:articleSlug", middleware.ValidateURIMiddleware(&protocol.ArticleURI{}))
	{
		articleSlugRouter.GET("", article.GetArticleInfoHandler)
		articleSlugRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateArticleBody{}), article.UpdateArticleHandler)
		articleSlugRouter.DELETE("", article.DeleteArticleHandler)
		articleSlugRouter.PUT("/status", middleware.ValidateBodyMiddleware(&protocol.UpdateArticleStatusBody{}), article.UpdateArticleStatusHandler)

		initArticleVersionRouter(articleSlugRouter)
		initArticleCommentRouter(articleSlugRouter)
	}
}

func initUserTagRouter(r *gin.RouterGroup) {
	r.GET("tags", middleware.ValidateParamMiddleware(&protocol.PageParam{}), tag.ListUserTagsHandler)
	tagRouter := r.Group("/tag")
	{
		tagRouter.GET("", middleware.ValidateParamMiddleware(&protocol.QueryParam{}), tag.QueryUserTagHandler)
	}
}

func initUserCategoryRouter(r *gin.RouterGroup) {
	r.GET("rootCategory", category.ListRootCategoriesHandler)
	categoryRouter := r.Group("/category")
	{
		categoryRouter.POST("", middleware.ValidateBodyMiddleware(&protocol.CreateCategoryBody{}), category.CreateCategoryHandler)
	}

	categoryIDRouter := categoryRouter.Group("/:categoryID", middleware.ValidateURIMiddleware(&protocol.CategoryURI{}))
	{
		categoryIDRouter.GET("", category.GetCategoryInfoHandler)
		categoryIDRouter.DELETE("", category.DeleteCategoryHandler)
		categoryIDRouter.PUT("", middleware.ValidateBodyMiddleware(&protocol.UpdateCategoryBody{}), category.UpdateCategoryInfoHandler)
		categoryIDRouter.GET("subCategories", middleware.ValidateParamMiddleware(&protocol.PageParam{}), category.ListChildrenCategoriesHandler)
		categoryIDRouter.GET("subArticles", middleware.ValidateParamMiddleware(&protocol.PageParam{}), category.ListChildrenArticlesHandler)
	}
}

func initUserOperationRouter(r *gin.RouterGroup) {
	operationRouter := r.Group("/operation")
	{
		initUserLikeRouter(operationRouter)
	}
}

func initUserLikeRouter(r *gin.RouterGroup) {
	userLikeRouter := r.Group("/like")
	{
		userLikeRouter.POST("article", middleware.ValidateBodyMiddleware(&protocol.LikeArticleBody{}), like.UserLikeArticleHandler)
		userLikeRouter.POST("comment", middleware.ValidateBodyMiddleware(&protocol.LikeCommentBody{}), like.UserLikeCommentHandler)
		userLikeRouter.POST("tag", middleware.ValidateBodyMiddleware(&protocol.LikeTagBody{}), like.UserLikeTagHandler)
	}
}

func initArticleVersionRouter(r *gin.RouterGroup) {
	r.GET("versions", middleware.ValidateParamMiddleware(&protocol.PageParam{}), article_version.ListArticleVersionsHandler)
	articleVersionRouter := r.Group("/version")
	{
		articleVersionRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "userID", protocol.CodeCreateArticleVersionRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleVersionBody{}),
			article_version.CreateArticleVersionHandler,
		)
		articleVersionNumberRouter := articleVersionRouter.Group("/v:version", middleware.ValidateURIMiddleware(&protocol.ArticleVersionURI{}))
		{
			articleVersionNumberRouter.GET("", article_version.GetArticleVersionInfoHandler)
		}
	}
}

func initArticleCommentRouter(r *gin.RouterGroup) {
	r.GET("comments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), comment.ListArticleCommentsHandler)
	commentRouter := r.Group("/comment")
	{
		commentRouter.POST(
			"",
			middleware.RateLimiterMiddleware(10*time.Second, 1, "userID", protocol.CodeCreateCommentRateLimitError),
			middleware.ValidateBodyMiddleware(&protocol.CreateArticleCommentBody{}),
			comment.CreateArticleCommentHandler,
		)
		commentIDRouter := commentRouter.Group("/:commentID", middleware.ValidateURIMiddleware(&protocol.CommentURI{}))
		{
			commentIDRouter.GET("", comment.GetCommentInfoHandler)
			commentIDRouter.DELETE("", comment.DeleteCommentHandler)
			commentIDRouter.GET("subComments", middleware.ValidateParamMiddleware(&protocol.PageParam{}), comment.ListChildrenCommentsHandler)
		}
	}
}
