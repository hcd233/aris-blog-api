package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HTTPResponse[*dto.CreateArticleResponse], error)
	HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleRequest) (*protocol.HTTPResponse[*dto.GetArticleResponse], error)
	HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleBySlugRequest) (*protocol.HTTPResponse[*dto.GetArticleBySlugResponse], error)
	HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
	HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
	HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
	HandleListArticles(ctx context.Context, req *dto.ListArticleRequest) (*protocol.HTTPResponse[*dto.ListArticleResponse], error)
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

func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HTTPResponse[*dto.CreateArticleResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateArticle(ctx, req))
}

func (h *articleHandler) HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleRequest) (*protocol.HTTPResponse[*dto.GetArticleResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetArticleInfo(ctx, req))
}

func (h *articleHandler) HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleBySlugRequest) (*protocol.HTTPResponse[*dto.GetArticleBySlugResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetArticleInfoBySlug(ctx, req))
}

func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateArticle(ctx, req))
}

func (h *articleHandler) HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateArticleStatus(ctx, req))
}

func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteArticle(ctx, req))
}

func (h *articleHandler) HandleListArticles(ctx context.Context, req *dto.ListArticleRequest) (*protocol.HTTPResponse[*dto.ListArticleResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListArticles(ctx, req))
}
