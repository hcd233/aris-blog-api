package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	humahandler "github.com/hcd233/aris-blog-api/internal/huma/handler"
)

// RegisterHumaRoutes 使用 Huma + net/http 适配器，并通过 Fiber adaptor 挂载到 Fiber
func RegisterHumaRoutes(app *fiber.App) {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Aris Blog API", "1.0.0"))

	// 实例化各模块 handler
	article := humahandler.NewArticleHandlers()
	category := humahandler.NewCategoryHandlers()
	tag := humahandler.NewTagHandlers()
	comment := humahandler.NewCommentHandlers()
	av := humahandler.NewArticleVersionHandlers()
	user := humahandler.NewUserHandlers()
	oauthGithub := humahandler.NewGithubOauth2Handlers()
	ai := humahandler.NewAIHandlers()
	token := humahandler.NewTokenHandlers()
	asset := humahandler.NewAssetHandlers()

	// 文章
	huma.Post(api, "/v1/article", article.HandleCreateArticle)
	huma.Get(api, "/v1/article/{articleID}", article.HandleGetArticleInfo)
	huma.Get(api, "/v1/article/slug/{authorName}/{articleSlug}", article.HandleGetArticleInfoBySlug)
	huma.Patch(api, "/v1/article/{articleID}", article.HandleUpdateArticle)
	huma.Put(api, "/v1/article/{articleID}/status", article.HandleUpdateArticleStatus)
	huma.Delete(api, "/v1/article/{articleID}", article.HandleDeleteArticle)
	huma.Get(api, "/v1/article/list", article.HandleListArticles)

	// 分类
	huma.Post(api, "/v1/category", category.HandleCreateCategory)
	huma.Get(api, "/v1/category/{categoryID}", category.HandleGetCategoryInfo)
	huma.Get(api, "/v1/category/root", category.HandleGetRootCategories)
	huma.Patch(api, "/v1/category/{categoryID}", category.HandleUpdateCategoryInfo)
	huma.Delete(api, "/v1/category/{categoryID}", category.HandleDeleteCategory)
	huma.Get(api, "/v1/category/{categoryID}/subCategories", category.HandleListChildrenCategories)
	huma.Get(api, "/v1/category/{categoryID}/subArticles", category.HandleListChildrenArticles)

	// 标签
	huma.Post(api, "/v1/tag", tag.HandleCreateTag)
	huma.Get(api, "/v1/tag/{tagID}", tag.HandleGetTagInfo)
	huma.Patch(api, "/v1/tag/{tagID}", tag.HandleUpdateTag)
	huma.Delete(api, "/v1/tag/{tagID}", tag.HandleDeleteTag)
	huma.Get(api, "/v1/tag/list", tag.HandleListTags)

	// 评论
	huma.Post(api, "/v1/comment", comment.HandleCreateArticleComment)
	huma.Delete(api, "/v1/comment/{commentID}", comment.HandleDeleteComment)
	huma.Get(api, "/v1/comment/article/{articleID}/list", comment.HandleListArticleComments)
	huma.Get(api, "/v1/comment/{commentID}/subComments", comment.HandleListChildrenComments)

	// 文章版本
	huma.Post(api, "/v1/article/{articleID}/version", av.HandleCreateArticleVersion)
	huma.Get(api, "/v1/article/{articleID}/version/v{version}", av.HandleGetArticleVersionInfo)
	huma.Get(api, "/v1/article/{articleID}/version/latest", av.HandleGetLatestArticleVersionInfo)
	huma.Get(api, "/v1/article/{articleID}/version/list", av.HandleListArticleVersions)

	// 用户
	huma.Get(api, "/v1/user/current", user.HandleGetCurUserInfo)
	huma.Get(api, "/v1/user/{userID}", user.HandleGetUserInfo)
	huma.Patch(api, "/v1/user", user.HandleUpdateInfo)

	// OAuth2（示例：仅注册 Github 登录/回调）
	huma.Get(api, "/v1/oauth2/github/login", oauthGithub.HandleLogin)
	huma.Get(api, "/v1/oauth2/github/callback", oauthGithub.HandleCallback)

	// AI
	huma.Get(api, "/v1/ai/prompt/{taskName}/v{version}", ai.HandleGetPrompt)
	huma.Get(api, "/v1/ai/prompt/{taskName}/latest", ai.HandleGetLatestPrompt)
	huma.Get(api, "/v1/ai/prompt/{taskName}", ai.HandleListPrompt)
	huma.Post(api, "/v1/ai/prompt/{taskName}", ai.HandleCreatePrompt)
	huma.Post(api, "/v1/ai/app/creator/contentCompletion", ai.HandleGenerateContentCompletion)
	huma.Post(api, "/v1/ai/app/creator/articleSummary", ai.HandleGenerateArticleSummary)
	huma.Post(api, "/v1/ai/app/creator/articleTranslation", ai.HandleGenerateArticleTranslation)
	huma.Post(api, "/v1/ai/app/reader/articleQA", ai.HandleGenerateArticleQA)
	huma.Post(api, "/v1/ai/app/reader/termExplaination", ai.HandleGenerateTermExplaination)

	// Token
	huma.Post(api, "/v1/token/refresh", token.HandleRefreshToken)

	// 资源
	huma.Get(api, "/v1/asset/like/articles", asset.HandleListUserLikeArticles)
	huma.Get(api, "/v1/asset/like/comments", asset.HandleListUserLikeComments)
	huma.Get(api, "/v1/asset/like/tags", asset.HandleListUserLikeTags)
	huma.Get(api, "/v1/asset/object/images", asset.HandleListImages)
	huma.Get(api, "/v1/asset/object/image/{objectName}", asset.HandleGetImage)
	huma.Delete(api, "/v1/asset/object/image/{objectName}", asset.HandleDeleteImage)
	huma.Get(api, "/v1/asset/view/articles", asset.HandleListUserViewArticles)
	huma.Delete(api, "/v1/asset/view/article/{viewID}", asset.HandleDeleteUserView)

	// 将 net/http mux 挂载到 Fiber
	app.All("/docs", adaptor.HTTPHandler(mux))
	app.All("/openapi*", adaptor.HTTPHandler(mux))
	app.All("/v1/*", adaptor.HTTPHandler(mux))
}
