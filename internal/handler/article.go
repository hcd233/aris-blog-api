package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
//
//	author centonhuang
//	update 2025-10-31 04:50:00
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *ArticleCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleResponse], error)
	HandleGetArticleInfo(ctx context.Context, req *ArticleGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleInfoResponse], error)
	HandleGetArticleInfoBySlug(ctx context.Context, req *ArticleGetBySlugRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleInfoBySlugResponse], error)
	HandleUpdateArticle(ctx context.Context, req *ArticleUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateArticleResponse], error)
	HandleUpdateArticleStatus(ctx context.Context, req *ArticleUpdateStatusRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateArticleStatusResponse], error)
	HandleDeleteArticle(ctx context.Context, req *ArticleDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteArticleResponse], error)
	HandleListArticles(ctx context.Context, req *ArticleListRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticlesResponse], error)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
//
//	return ArticleHandler
//	author centonhuang
//	update 2025-10-31 04:50:00
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

// ArticlePathParam 文章路径参数
type ArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"文章 ID"`
}

// ArticleSlugPathParam 文章别名路径参数
type ArticleSlugPathParam struct {
	AuthorName  string `path:"authorName" doc:"作者名称"`
	ArticleSlug string `path:"articleSlug" doc:"文章别名"`
}

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Body *protocol.CreateArticleBody `json:"body" doc:"创建文章请求体"`
}

// ArticleGetRequest 获取文章请求
type ArticleGetRequest struct {
	ArticlePathParam
}

// ArticleGetBySlugRequest 按别名获取文章请求
type ArticleGetBySlugRequest struct {
	ArticleSlugPathParam
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	ArticlePathParam
	Body *protocol.UpdateArticleBody `json:"body" doc:"更新文章请求体"`
}

// ArticleUpdateStatusRequest 更新文章状态请求
type ArticleUpdateStatusRequest struct {
	ArticlePathParam
	Body *protocol.UpdateArticleStatusBody `json:"body" doc:"更新文章状态请求体"`
}

// ArticleDeleteRequest 删除文章请求
type ArticleDeleteRequest struct {
	ArticlePathParam
}

// ArticleListRequest 列出文章请求
type ArticleListRequest struct {
	PaginationQuery
}

// HandleCreateArticle 创建文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleCreateRequest
//	return *protocol.HumaHTTPResponse[*protocol.CreateArticleResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *ArticleCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.CreateArticleRequest{
		UserID:     userID,
		Title:      req.Body.Title,
		Slug:       req.Body.Slug,
		CategoryID: req.Body.CategoryID,
		Tags:       req.Body.Tags,
	}

	return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, serviceReq))
}

// HandleGetArticleInfo 获取文章详情
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleGetRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetArticleInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleGetArticleInfo(ctx context.Context, req *ArticleGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleInfoResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetArticleInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	return util.WrapHTTPResponse(h.svc.GetArticleInfo(ctx, serviceReq))
}

// HandleGetArticleInfoBySlug 通过别名获取文章详情
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleGetBySlugRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetArticleInfoBySlugResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleGetArticleInfoBySlug(ctx context.Context, req *ArticleGetBySlugRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleInfoBySlugResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetArticleInfoBySlugRequest{
		UserID:      userID,
		AuthorName:  req.AuthorName,
		ArticleSlug: req.ArticleSlug,
	}

	return util.WrapHTTPResponse(h.svc.GetArticleInfoBySlug(ctx, serviceReq))
}

// HandleUpdateArticle 更新文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleUpdateRequest
//	return *protocol.HumaHTTPResponse[*protocol.UpdateArticleResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *ArticleUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateArticleResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.UpdateArticleRequest{
		UserID:            userID,
		ArticleID:         req.ArticleID,
		UpdatedTitle:      req.Body.Title,
		UpdatedSlug:       req.Body.Slug,
		UpdatedCategoryID: req.Body.CategoryID,
	}

	return util.WrapHTTPResponse(h.svc.UpdateArticle(ctx, serviceReq))
}

// HandleUpdateArticleStatus 更新文章状态
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleUpdateStatusRequest
//	return *protocol.HumaHTTPResponse[*protocol.UpdateArticleStatusResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleUpdateArticleStatus(ctx context.Context, req *ArticleUpdateStatusRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateArticleStatusResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.UpdateArticleStatusRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Status:    req.Body.Status,
	}

	return util.WrapHTTPResponse(h.svc.UpdateArticleStatus(ctx, serviceReq))
}

// HandleDeleteArticle 删除文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleDeleteRequest
//	return *protocol.HumaHTTPResponse[*protocol.DeleteArticleResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *ArticleDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteArticleResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.DeleteArticleRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	return util.WrapHTTPResponse(h.svc.DeleteArticle(ctx, serviceReq))
}

// HandleListArticles 列出文章
//
//	receiver h *articleHandler
//	param ctx context.Context
//	param req *ArticleListRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListArticlesResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:50:00
func (h *articleHandler) HandleListArticles(ctx context.Context, req *ArticleListRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticlesResponse], error) {
	serviceReq := &protocol.ListArticlesRequest{
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListArticles(ctx, serviceReq))
}
