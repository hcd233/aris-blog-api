package router

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

// RegisterHumaRouter 注册 Huma 路由
func RegisterHumaRouter(app *fiber.App) {
	// 创建 Huma API 实例
	config := huma.DefaultConfig("Aris-blog API", "1.0.0")
	api := humafiber.New(app, config)

	// Ping 路由
	pingHandler := func(ctx context.Context, input *struct{}) (*protocol.PingResponse, error) {
		return &protocol.PingResponse{Message: "ok"}, nil
	}
	huma.Get(api, "/", pingHandler)

	// v1 路由组
	v1 := huma.NewGroup(api, "/v1")

	// 用户路由
	initHumaUserRouter(v1)

	// 文章路由
	initHumaArticleRouter(v1)

	// 标签路由
	initHumaTagRouter(v1)

	// 分类路由
	initHumaCategoryRouter(v1)

	// 评论路由
	initHumaCommentRouter(v1)
}

// initHumaUserRouter 初始化用户路由
func initHumaUserRouter(v1 *huma.Group) {
	// 当前用户信息
	huma.Get(v1, "/user/current", handler.GetCurUserInfoHuma)

	// 用户信息
	huma.Get(v1, "/user/{userID}", handler.GetUserInfoHuma)

	// 更新用户信息
	huma.Patch(v1, "/user", handler.UpdateUserInfoHuma)
}

// initHumaArticleRouter 初始化文章路由
func initHumaArticleRouter(v1 *huma.Group) {
	// 文章列表
	huma.Get(v1, "/article/list", handler.ListArticlesHuma)

	// 通过 slug 获取文章
	huma.Get(v1, "/article/slug/{authorName}/{articleSlug}", handler.GetArticleInfoBySlugHuma)

	// 创建文章
	huma.Post(v1, "/article", handler.CreateArticleHuma)

	// 文章详情
	huma.Get(v1, "/article/{articleID}", handler.GetArticleInfoHuma)

	// 更新文章
	huma.Patch(v1, "/article/{articleID}", handler.UpdateArticleHuma)

	// 删除文章
	huma.Delete(v1, "/article/{articleID}", handler.DeleteArticleHuma)

	// 更新文章状态
	huma.Put(v1, "/article/{articleID}/status", handler.UpdateArticleStatusHuma)
}

// initHumaTagRouter 初始化标签路由
func initHumaTagRouter(v1 *huma.Group) {
	// 标签列表
	huma.Get(v1, "/tag/list", handler.ListTagsHuma)

	// 创建标签
	huma.Post(v1, "/tag", handler.CreateTagHuma)

	// 标签详情
	huma.Get(v1, "/tag/{tagID}", handler.GetTagInfoHuma)

	// 更新标签
	huma.Patch(v1, "/tag/{tagID}", handler.UpdateTagHuma)

	// 删除标签
	huma.Delete(v1, "/tag/{tagID}", handler.DeleteTagHuma)
}

// initHumaCategoryRouter 初始化分类路由
func initHumaCategoryRouter(v1 *huma.Group) {
	// 获取根分类
	huma.Get(v1, "/category/root", handler.GetRootCategoriesHuma)

	// 创建分类
	huma.Post(v1, "/category", handler.CreateCategoryHuma)

	// 分类详情
	huma.Get(v1, "/category/{categoryID}", handler.GetCategoryInfoHuma)

	// 更新分类
	huma.Patch(v1, "/category/{categoryID}", handler.UpdateCategoryHuma)

	// 删除分类
	huma.Delete(v1, "/category/{categoryID}", handler.DeleteCategoryHuma)

	// 列出子分类
	huma.Get(v1, "/category/{categoryID}/subCategories", handler.ListChildrenCategoriesHuma)
}

// initHumaCommentRouter 初始化评论路由
func initHumaCommentRouter(v1 *huma.Group) {
	// 创建评论
	huma.Post(v1, "/comment", handler.CreateCommentHuma)

	// 列出文章评论
	huma.Get(v1, "/comment/article/{articleID}/list", handler.ListArticleCommentsHuma)

	// 删除评论
	huma.Delete(v1, "/comment/{commentID}", handler.DeleteCommentHuma)

	// 列出子评论
	huma.Get(v1, "/comment/{commentID}/subComments", handler.ListChildrenCommentsHuma)
}
