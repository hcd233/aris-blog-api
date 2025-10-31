package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleVersionHandler 文章版本处理器
//
//	author centonhuang
//	update 2025-10-30
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(ctx context.Context, req *dto.CreateArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleVersionResponse], error)
	HandleGetArticleVersionInfo(ctx context.Context, req *dto.GetArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleVersionInfoResponse], error)
	HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.GetLatestArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestArticleVersionInfoResponse], error)
	HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionsResponse], error)
}

type articleVersionHandler struct {
	svc service.ArticleVersionService
}

// NewArticleVersionHandler 创建文章版本处理器
//
//	return ArticleVersionHandler
//	author centonhuang
//	update 2025-10-30
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		svc: service.NewArticleVersionService(),
	}
}

func (h *articleVersionHandler) HandleCreateArticleVersion(ctx context.Context, req *dto.CreateArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleVersionResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.CreateArticleVersionRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Content:   req.Body.Content,
	}

	svcRsp, err := h.svc.CreateArticleVersion(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.CreateArticleVersionResponse](nil, err)
	}

	rsp := &dto.CreateArticleVersionResponse{
		ArticleVersion: &dto.ArticleVersion{
			ArticleVersionID: svcRsp.ArticleVersion.ArticleVersionID,
			ArticleID:        svcRsp.ArticleVersion.ArticleID,
			VersionID:        svcRsp.ArticleVersion.VersionID,
			Content:          svcRsp.ArticleVersion.Content,
			CreatedAt:        svcRsp.ArticleVersion.CreatedAt,
			UpdatedAt:        svcRsp.ArticleVersion.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleVersionHandler) HandleGetArticleVersionInfo(ctx context.Context, req *dto.GetArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleVersionInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		VersionID: req.Version,
	}

	svcRsp, err := h.svc.GetArticleVersionInfo(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetArticleVersionInfoResponse](nil, err)
	}

	rsp := &dto.GetArticleVersionInfoResponse{
		Version: &dto.ArticleVersion{
			ArticleVersionID: svcRsp.Version.ArticleVersionID,
			ArticleID:        svcRsp.Version.ArticleID,
			VersionID:        svcRsp.Version.VersionID,
			Content:          svcRsp.Version.Content,
			CreatedAt:        svcRsp.Version.CreatedAt,
			UpdatedAt:        svcRsp.Version.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.GetLatestArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestArticleVersionInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetLatestArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	svcRsp, err := h.svc.GetLatestArticleVersionInfo(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetLatestArticleVersionInfoResponse](nil, err)
	}

	rsp := &dto.GetLatestArticleVersionInfoResponse{
		Version: &dto.ArticleVersion{
			ArticleVersionID: svcRsp.Version.ArticleVersionID,
			ArticleID:        svcRsp.Version.ArticleID,
			VersionID:        svcRsp.Version.VersionID,
			Content:          svcRsp.Version.Content,
			CreatedAt:        svcRsp.Version.CreatedAt,
			UpdatedAt:        svcRsp.Version.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleVersionHandler) HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionsResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListArticleVersionsRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListArticleVersions(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListArticleVersionsResponse](nil, err)
	}

	versions := make([]*dto.ArticleVersion, len(svcRsp.Versions))
	for i, version := range svcRsp.Versions {
		versions[i] = &dto.ArticleVersion{
			ArticleVersionID: version.ArticleVersionID,
			ArticleID:        version.ArticleID,
			VersionID:        version.VersionID,
			Content:          version.Content,
			CreatedAt:        version.CreatedAt,
			UpdatedAt:        version.UpdatedAt,
		}
	}

	rsp := &dto.ListArticleVersionsResponse{
		Versions: versions,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}
