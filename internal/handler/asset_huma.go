package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// AssetHandlerForHuma 资产处理器（Huma版本）
type AssetHandlerForHuma interface {
	HandleListUserLikeArticles(ctx context.Context, req *dto.ListUserLikeArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeArticlesResponse], error)
	HandleListUserLikeComments(ctx context.Context, req *dto.ListUserLikeCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeCommentsResponse], error)
	HandleListUserLikeTags(ctx context.Context, req *dto.ListUserLikeTagsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeTagsResponse], error)
	HandleListImages(ctx context.Context, req *dto.ListImagesRequest) (*protocol.HumaHTTPResponse[*dto.ListImagesResponse], error)
	HandleUploadImage(ctx context.Context, req *dto.UploadImageRequest) (*protocol.HumaHTTPResponse[*dto.UploadImageResponse], error)
	HandleGetImage(ctx context.Context, req *dto.GetImageRequest) (*protocol.HumaHTTPResponse[*dto.GetImageResponse], error)
	HandleDeleteImage(ctx context.Context, req *dto.DeleteImageRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleListUserViewArticles(ctx context.Context, req *dto.ListUserViewArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserViewArticlesResponse], error)
	HandleDeleteUserView(ctx context.Context, req *dto.DeleteUserViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
}

type assetHandlerForHuma struct {
	svc service.AssetService
}

// NewAssetHandlerForHuma 创建资产处理器（Huma版本）
func NewAssetHandlerForHuma() AssetHandlerForHuma {
	return &assetHandlerForHuma{
		svc: service.NewAssetService(),
	}
}

func (h *assetHandlerForHuma) HandleListUserLikeArticles(ctx context.Context, req *dto.ListUserLikeArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeArticlesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeArticles(ctx, req))
}

func (h *assetHandlerForHuma) HandleListUserLikeComments(ctx context.Context, req *dto.ListUserLikeCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeCommentsResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeComments(ctx, req))
}

func (h *assetHandlerForHuma) HandleListUserLikeTags(ctx context.Context, req *dto.ListUserLikeTagsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeTagsResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeTags(ctx, req))
}

func (h *assetHandlerForHuma) HandleListImages(ctx context.Context, req *dto.ListImagesRequest) (*protocol.HumaHTTPResponse[*dto.ListImagesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListImages(ctx, req))
}

func (h *assetHandlerForHuma) HandleUploadImage(ctx context.Context, req *dto.UploadImageRequest) (*protocol.HumaHTTPResponse[*dto.UploadImageResponse], error) {
	// 文件上传需要特殊处理
	// 由于huma的限制，我们需要从fiber.Ctx中获取文件
	// 这里暂时返回未实现错误
	logger.WithCtx(ctx).Warn("[AssetHandler] UploadImage not yet implemented for Huma")
	return nil, protocol.ErrNoImplement
}

func (h *assetHandlerForHuma) HandleGetImage(ctx context.Context, req *dto.GetImageRequest) (*protocol.HumaHTTPResponse[*dto.GetImageResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	
	internalReq := &dto.InternalGetImageRequest{
		UserID:    userID,
		ImageName: req.ObjectName,
		Quality:   req.Quality,
	}
	
	rsp, err := h.svc.GetImage(ctx, internalReq)
	if err != nil {
		return util.WrapHTTPResponse(rsp, err)
	}
	
	// 对于图片获取，我们返回presigned URL，客户端需要重定向
	return util.WrapHTTPResponse(rsp, err)
}

func (h *assetHandlerForHuma) HandleDeleteImage(ctx context.Context, req *dto.DeleteImageRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)
	
	internalReq := &dto.InternalDeleteImageRequest{
		UserID:    userID,
		ImageName: req.ObjectName,
	}
	
	return util.WrapHTTPResponse(h.svc.DeleteImage(ctx, internalReq))
}

func (h *assetHandlerForHuma) HandleListUserViewArticles(ctx context.Context, req *dto.ListUserViewArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserViewArticlesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserViewArticles(ctx, req))
}

func (h *assetHandlerForHuma) HandleDeleteUserView(ctx context.Context, req *dto.DeleteUserViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteUserView(ctx, req))
}
