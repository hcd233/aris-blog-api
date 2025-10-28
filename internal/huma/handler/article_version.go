package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type ArticleVersionHandlers struct{ svc service.ArticleVersionService }

func NewArticleVersionHandlers() *ArticleVersionHandlers {
	return &ArticleVersionHandlers{svc: service.NewArticleVersionService()}
}

type (
	createVersionInput struct {
		authHeader
		humadto.ArticlePath
		humadto.CreateArticleVersionInput
	}
	versionPathInput struct {
		authHeader
		humadto.ArticleVersionPath
	}
	articlePathInputAV struct {
		authHeader
		humadto.ArticlePath
	}
	listVersionsInput struct {
		authHeader
		humadto.ArticlePath
		humadto.PaginateParam
	}
)

func (h *ArticleVersionHandlers) HandleCreateArticleVersion(ctx context.Context, input *createVersionInput) (*humadto.Output[protocol.CreateArticleVersionResponse], error) {
	req := &protocol.CreateArticleVersionRequest{UserID: input.UserID, ArticleID: input.ArticleID, Content: input.Body.Content}
	rsp, err := h.svc.CreateArticleVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreateArticleVersionResponse]{Body: *rsp}, nil
}

func (h *ArticleVersionHandlers) HandleGetArticleVersionInfo(ctx context.Context, input *versionPathInput) (*humadto.Output[protocol.GetArticleVersionInfoResponse], error) {
	req := &protocol.GetArticleVersionInfoRequest{UserID: input.UserID, ArticleID: input.ArticleID, VersionID: input.Version}
	rsp, err := h.svc.GetArticleVersionInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetArticleVersionInfoResponse]{Body: *rsp}, nil
}

func (h *ArticleVersionHandlers) HandleGetLatestArticleVersionInfo(ctx context.Context, input *articlePathInputAV) (*humadto.Output[protocol.GetLatestArticleVersionInfoResponse], error) {
	req := &protocol.GetLatestArticleVersionInfoRequest{UserID: input.UserID, ArticleID: input.ArticleID}
	rsp, err := h.svc.GetLatestArticleVersionInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetLatestArticleVersionInfoResponse]{Body: *rsp}, nil
}

func (h *ArticleVersionHandlers) HandleListArticleVersions(ctx context.Context, input *listVersionsInput) (*humadto.Output[protocol.ListArticleVersionsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListArticleVersionsRequest{UserID: input.UserID, ArticleID: input.ArticleID, PaginateParam: p}
	rsp, err := h.svc.ListArticleVersions(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListArticleVersionsResponse]{Body: *rsp}, nil
}
