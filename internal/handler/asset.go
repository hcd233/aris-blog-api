package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// AssetHandler 资产处理器
type AssetHandler interface {
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

type assetHandler struct {
	svc service.AssetService
}

// NewAssetHandler 创建资产处理器
func NewAssetHandler() AssetHandler {
	return &assetHandler{
		svc: service.NewAssetService(),
	}
}

func (h *assetHandler) HandleListUserLikeArticles(ctx context.Context, req *dto.ListUserLikeArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeArticlesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeArticles(ctx, req))
}

func (h *assetHandler) HandleListUserLikeComments(ctx context.Context, req *dto.ListUserLikeCommentsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeCommentsResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeComments(ctx, req))
}

func (h *assetHandler) HandleListUserLikeTags(ctx context.Context, req *dto.ListUserLikeTagsRequest) (*protocol.HumaHTTPResponse[*dto.ListUserLikeTagsResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserLikeTags(ctx, req))
}

func (h *assetHandler) HandleListImages(ctx context.Context, req *dto.ListImagesRequest) (*protocol.HumaHTTPResponse[*dto.ListImagesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListImages(ctx, req))
}

func (h *assetHandler) HandleUploadImage(ctx context.Context, req *dto.UploadImageRequest) (*protocol.HumaHTTPResponse[*dto.UploadImageResponse], error) {
	return util.WrapHTTPResponse(h.svc.UploadImage(ctx, req))
}

func (h *assetHandler) HandleGetImage(ctx context.Context, req *dto.GetImageRequest) (*protocol.HumaHTTPResponse[*dto.GetImageResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetImage(ctx, req))
}

func (h *assetHandler) HandleDeleteImage(ctx context.Context, req *dto.DeleteImageRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteImage(ctx, req))
}

func (h *assetHandler) HandleListUserViewArticles(ctx context.Context, req *dto.ListUserViewArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListUserViewArticlesResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListUserViewArticles(ctx, req))
}

func (h *assetHandler) HandleDeleteUserView(ctx context.Context, req *dto.DeleteUserViewRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteUserView(ctx, req))
}
