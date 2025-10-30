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
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		svc: service.NewArticleVersionService(),
	}
}

// HandleCreateArticleVersion 创建文章版本
func (h *articleVersionHandler) HandleCreateArticleVersion(ctx context.Context, req *dto.CreateArticleVersionRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleVersionResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.CreateArticleVersionRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Content:   req.Body.Content,
	}

	serviceRsp, err := h.svc.CreateArticleVersion(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.CreateArticleVersionResponse{
		ArticleVersion: convertArticleVersion(serviceRsp.ArticleVersion),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleGetArticleVersionInfo 获取文章版本信息
func (h *articleVersionHandler) HandleGetArticleVersionInfo(ctx context.Context, req *dto.GetArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleVersionInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.GetArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		VersionID: req.Version,
	}

	serviceRsp, err := h.svc.GetArticleVersionInfo(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.GetArticleVersionInfoResponse{
		Version: convertArticleVersion(serviceRsp.Version),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleGetLatestArticleVersionInfo 获取最新文章版本信息
func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(ctx context.Context, req *dto.GetLatestArticleVersionInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetLatestArticleVersionInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.GetLatestArticleVersionInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	serviceRsp, err := h.svc.GetLatestArticleVersionInfo(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	rsp := &dto.GetLatestArticleVersionInfoResponse{
		Version: convertArticleVersion(serviceRsp.Version),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleListArticleVersions 列出文章版本
func (h *articleVersionHandler) HandleListArticleVersions(ctx context.Context, req *dto.ListArticleVersionsRequest) (*protocol.HumaHTTPResponse[*dto.ListArticleVersionsResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	serviceReq := &protocol.ListArticleVersionsRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     req.Page,
				PageSize: req.PageSize,
			},
		},
	}

	serviceRsp, err := h.svc.ListArticleVersions(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	versions := make([]*dto.ArticleVersion, len(serviceRsp.Versions))
	for i, version := range serviceRsp.Versions {
		versions[i] = convertArticleVersion(version)
	}

	rsp := &dto.ListArticleVersionsResponse{
		Versions: versions,
		PageInfo: &dto.PageInfo{
			Page:     serviceRsp.PageInfo.Page,
			PageSize: serviceRsp.PageInfo.PageSize,
			Total:    serviceRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// convertArticleVersion 转换文章版本模型
func convertArticleVersion(version *protocol.ArticleVersion) *dto.ArticleVersion {
	if version == nil {
		return nil
	}

	return &dto.ArticleVersion{
		ArticleVersionID: version.ArticleVersionID,
		ArticleID:        version.ArticleID,
		VersionID:        version.VersionID,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt,
		UpdatedAt:        version.UpdatedAt,
	}
}
