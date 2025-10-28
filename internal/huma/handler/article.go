package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// ArticleHandlers 封装文章相关的 Huma 风格处理函数（纯函数形式）
// 不直接依赖路由，后续由适配层注册到具体 Huma API
type ArticleHandlers struct {
	svc service.ArticleService
}

func NewArticleHandlers() *ArticleHandlers {
	return &ArticleHandlers{svc: service.NewArticleService()}
}

// 为简化演示，先通过 Header 注入用户 ID；后续适配 JWT 中间件
type authHeader struct {
	UserID uint `header:"X-User-ID" doc:"Temporary user ID from header; will be replaced by JWT"`
}

// 输入聚合：Body、Path、Query、Header
type (
	createArticleInput struct {
		authHeader
		humadto.CreateArticleInput
	}

	articlePathInput struct {
		authHeader
		humadto.ArticlePath
	}

	articleSlugPathInput struct {
		authHeader
		humadto.ArticleSlugPath
	}

	updateArticleInput struct {
		authHeader
		humadto.ArticlePath
		humadto.UpdateArticleInput
	}

	updateArticleStatusInput struct {
		authHeader
		humadto.ArticlePath
		humadto.UpdateArticleStatusInput
	}

	listArticlesInput struct {
		humadto.PaginateParam
	}
)

func (h *ArticleHandlers) HandleCreateArticle(ctx context.Context, input *createArticleInput) (*humadto.Output[protocol.CreateArticleResponse], error) {
	req := &protocol.CreateArticleRequest{
		UserID:     input.UserID,
		Title:      input.Body.Title,
		Slug:       input.Body.Slug,
		CategoryID: input.Body.CategoryID,
		Tags:       input.Body.Tags,
	}
	rsp, err := h.svc.CreateArticle(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreateArticleResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleGetArticleInfo(ctx context.Context, input *articlePathInput) (*humadto.Output[protocol.GetArticleInfoResponse], error) {
	req := &protocol.GetArticleInfoRequest{
		UserID:    input.UserID,
		ArticleID: input.ArticleID,
	}
	rsp, err := h.svc.GetArticleInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetArticleInfoResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleGetArticleInfoBySlug(ctx context.Context, input *articleSlugPathInput) (*humadto.Output[protocol.GetArticleInfoBySlugResponse], error) {
	req := &protocol.GetArticleInfoBySlugRequest{
		UserID:      input.UserID,
		AuthorName:  input.AuthorName,
		ArticleSlug: input.ArticleSlug,
	}
	rsp, err := h.svc.GetArticleInfoBySlug(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetArticleInfoBySlugResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleUpdateArticle(ctx context.Context, input *updateArticleInput) (*humadto.Output[protocol.UpdateArticleResponse], error) {
	req := &protocol.UpdateArticleRequest{
		UserID:            input.UserID,
		ArticleID:         input.ArticleID,
		UpdatedTitle:      input.Body.Title,
		UpdatedSlug:       input.Body.Slug,
		UpdatedCategoryID: input.Body.CategoryID,
	}
	rsp, err := h.svc.UpdateArticle(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.UpdateArticleResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleUpdateArticleStatus(ctx context.Context, input *updateArticleStatusInput) (*humadto.Output[protocol.UpdateArticleStatusResponse], error) {
	req := &protocol.UpdateArticleStatusRequest{
		UserID:    input.UserID,
		ArticleID: input.ArticleID,
		Status:    input.Body.Status,
	}
	rsp, err := h.svc.UpdateArticleStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.UpdateArticleStatusResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleDeleteArticle(ctx context.Context, input *articlePathInput) (*humadto.Output[protocol.DeleteArticleResponse], error) {
	req := &protocol.DeleteArticleRequest{
		UserID:    input.UserID,
		ArticleID: input.ArticleID,
	}
	rsp, err := h.svc.DeleteArticle(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteArticleResponse]{Body: *rsp}, nil
}

func (h *ArticleHandlers) HandleListArticles(ctx context.Context, input *listArticlesInput) (*humadto.Output[protocol.ListArticlesResponse], error) {
	req := &protocol.ListArticlesRequest{PaginateParam: &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}}
	if input.PageParam != nil {
		req.PaginateParam.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		req.PaginateParam.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	rsp, err := h.svc.ListArticles(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListArticlesResponse]{Body: *rsp}, nil
}
