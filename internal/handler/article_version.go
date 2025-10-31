package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleVersionHandler 文章版本处理器
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(ctx context.Context, req *dto.ArticleVersionCreateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionCreateResponse], error)
	HandleGetArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionGetResponse], error)
	HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetLatestRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionGetLatestResponse], error)
	HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionResponse], error)
}

type articleVersionHandler struct {
	svc service.ArticleVersionService
}

// NewArticleVersionHandler 创建文章版本处理器
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		svc: service.NewArticleVersionService(),
	}
}

func (h *articleVersionHandler) HandleCreateArticleVersion(ctx context.Context, req *dto.ArticleVersionCreateRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionCreateResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateArticleVersion(ctx, req))
}

func (h *articleVersionHandler) HandleGetArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionGetResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetArticleVersionInfo(ctx, req))
}

func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetLatestRequest) (*protocol.HumaHTTPResponse[*dto.ArticleVersionGetLatestResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetLatestArticleVersionInfo(ctx, req))
}

func (h *articleVersionHandler) HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListArticleVersions(ctx, req))
}
