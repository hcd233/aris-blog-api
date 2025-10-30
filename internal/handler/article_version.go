package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleVersionHandler 文章版本处理器
//
//	author centonhuang
//	update 2025-10-31 05:05:00
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(ctx context.Context, req *ArticleVersionCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleVersionResponse], error)
	HandleGetArticleVersionInfo(ctx context.Context, req *ArticleVersionGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleVersionInfoResponse], error)
	HandleGetLatestArticleVersionInfo(ctx context.Context, req *ArticleVersionGetLatestRequest) (*protocol.HumaHTTPResponse[*protocol.GetLatestArticleVersionInfoResponse], error)
	HandleListArticleVersions(ctx context.Context, req *ArticleVersionListRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticleVersionsResponse], error)
}

type articleVersionHandler struct {
	svc service.ArticleVersionService
}

// NewArticleVersionHandler 创建文章版本处理器
//
//	return ArticleVersionHandler
//	author centonhuang
//	update 2025-10-31 05:05:00
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		svc: service.NewArticleVersionService(),
	}
}

// ArticleVersionArticlePathParam 文章路径参数
type ArticleVersionArticlePathParam struct {
	ArticleID uint `path:"articleID" doc:"文章 ID"`
}

// ArticleVersionPathParam 文章版本路径参数
type ArticleVersionPathParam struct {
	ArticleVersionArticlePathParam
	Version uint `path:"version" doc:"版本号"`
}

// ArticleVersionCreateRequest 创建文章版本请求
type ArticleVersionCreateRequest struct {
	ArticleVersionArticlePathParam
	Body *protocol.CreateArticleVersionBody `json:"body" doc:"创建文章版本请求体"`
}

// ArticleVersionGetRequest 获取文章版本详情请求
type ArticleVersionGetRequest struct {
	ArticleVersionPathParam
}

// ArticleVersionGetLatestRequest 获取最新文章版本请求
type ArticleVersionGetLatestRequest struct {
	ArticleVersionArticlePathParam
}

// ArticleVersionListRequest 列出文章版本请求
type ArticleVersionListRequest struct {
	ArticleVersionArticlePathParam
	PaginationQuery
}

// HandleCreateArticleVersion 创建文章版本
//
//	receiver h *articleVersionHandler
//	param ctx context.Context
//	param req *ArticleVersionCreateRequest
//	return *protocol.HumaHTTPResponse[*protocol.CreateArticleVersionResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:05:00
func (h *articleVersionHandler) HandleCreateArticleVersion(ctx context.Context, req *ArticleVersionCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateArticleVersionResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.CreateArticleVersionRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Content:   req.Body.Content,
	}

	return util.WrapHTTPResponse(h.svc.CreateArticleVersion(ctx, serviceReq))
}

// HandleGetArticleVersionInfo 获取文章版本信息
//
//	receiver h *articleVersionHandler
//	param ctx context.Context
//	param req *ArticleVersionGetRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetArticleVersionInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:05:00
func (h *articleVersionHandler) HandleGetArticleVersionInfo(ctx context.Context, req *ArticleVersionGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetArticleVersionInfoResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		VersionID: req.Version,
	}

	return util.WrapHTTPResponse(h.svc.GetArticleVersionInfo(ctx, serviceReq))
}

// HandleGetLatestArticleVersionInfo 获取最新文章版本
//
//	receiver h *articleVersionHandler
//	param ctx context.Context
//	param req *ArticleVersionGetLatestRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetLatestArticleVersionInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:05:00
func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(ctx context.Context, req *ArticleVersionGetLatestRequest) (*protocol.HumaHTTPResponse[*protocol.GetLatestArticleVersionInfoResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetLatestArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	return util.WrapHTTPResponse(h.svc.GetLatestArticleVersionInfo(ctx, serviceReq))
}

// HandleListArticleVersions 列出文章版本
//
//	receiver h *articleVersionHandler
//	param ctx context.Context
//	param req *ArticleVersionListRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListArticleVersionsResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 05:05:00
func (h *articleVersionHandler) HandleListArticleVersions(ctx context.Context, req *ArticleVersionListRequest) (*protocol.HumaHTTPResponse[*protocol.ListArticleVersionsResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.ListArticleVersionsRequest{
		UserID:        userID,
		ArticleID:     req.ArticleID,
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListArticleVersions(ctx, serviceReq))
}
