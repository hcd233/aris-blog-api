package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.ArticleCreateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleCreateResponse], error)
	HandleGetArticleInfo(ctx context.Context, req *dto.ArticleGetRequest) (*protocol.HumaHTTPResponse[*dto.ArticleGetResponse], error)
	HandleGetArticleInfoBySlug(ctx context.Context, req *dto.ArticleGetBySlugRequest) (*protocol.HumaHTTPResponse[*dto.ArticleGetBySlugResponse], error)
	HandleUpdateArticle(ctx context.Context, req *dto.ArticleUpdateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleUpdateResponse], error)
	HandleUpdateArticleStatus(ctx context.Context, req *dto.ArticleUpdateStatusRequest) (*protocol.HumaHTTPResponse[*dto.ArticleUpdateStatusResponse], error)
	HandleDeleteArticle(ctx context.Context, req *dto.ArticleDeleteRequest) (*protocol.HumaHTTPResponse[*dto.ArticleDeleteResponse], error)
	HandleListArticles(ctx context.Context, req *dto.ArticleListRequest) (*protocol.HumaHTTPResponse[*dto.ArticleListResponse], error)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.ArticleCreateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleCreateResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}

func (h *articleHandler) HandleGetArticleInfo(ctx context.Context, req *dto.ArticleGetRequest) (*protocol.HumaHTTPResponse[*dto.ArticleGetResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.GetArticleInfo(ctx, req))
}

func (h *articleHandler) HandleGetArticleInfoBySlug(ctx context.Context, req *dto.ArticleGetBySlugRequest) (*protocol.HumaHTTPResponse[*dto.ArticleGetBySlugResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.GetArticleInfoBySlug(ctx, req))
}

func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *dto.ArticleUpdateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleUpdateResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.UpdateArticle(ctx, req))
}

func (h *articleHandler) HandleUpdateArticleStatus(ctx context.Context, req *dto.ArticleUpdateStatusRequest) (*protocol.HumaHTTPResponse[*dto.ArticleUpdateStatusResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.UpdateArticleStatus(ctx, req))
}

func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *dto.ArticleDeleteRequest) (*protocol.HumaHTTPResponse[*dto.ArticleDeleteResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.DeleteArticle(ctx, req))
}

func (h *articleHandler) HandleListArticles(ctx context.Context, req *dto.ArticleListRequest) (*protocol.HumaHTTPResponse[*dto.ArticleListResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListArticles(ctx, req))
}
