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
	HandleCreateArticleVersion(ctx context.Context, req *dto.CreateArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleVersionResponse], error)
	HandleGetArticleVersionInfo(ctx context.Context, req *dto.GetArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleVersionResponse], error)
	HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.GetLatestArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestArticleVersionResponse], error)
	HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionsResponse], error)
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

func (h *articleVersionHandler) HandleCreateArticleVersion(ctx context.Context, req *dto.CreateArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleVersionResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateArticleVersion(ctx, req))
}

func (h *articleVersionHandler) HandleGetArticleVersionInfo(ctx context.Context, req *dto.GetArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleVersionResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetArticleVersionInfo(ctx, req))
}

func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.GetLatestArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestArticleVersionResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetLatestArticleVersionInfo(ctx, req))
}

func (h *articleVersionHandler) HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionsResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListArticleVersions(ctx, req))
}
