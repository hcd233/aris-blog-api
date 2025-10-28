package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type AssetHandlers struct{ svc service.AssetService }

func NewAssetHandlers() *AssetHandlers { return &AssetHandlers{svc: service.NewAssetService()} }

type (
	listWithPageInput struct {
		authHeader
		humadto.PaginateParam
	}
	objectPathInput struct {
		authHeader
		humadto.ObjectPath
	}
	objectWithImageParamInput struct {
		authHeader
		humadto.ObjectPath
		humadto.ImageParam
	}
)

func (h *AssetHandlers) HandleListUserLikeArticles(ctx context.Context, input *listWithPageInput) (*humadto.Output[protocol.ListUserLikeArticlesResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListUserLikeArticlesRequest{UserID: input.UserID, PaginateParam: p}
	rsp, err := h.svc.ListUserLikeArticles(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListUserLikeArticlesResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleListUserLikeComments(ctx context.Context, input *listWithPageInput) (*humadto.Output[protocol.ListUserLikeCommentsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListUserLikeCommentsRequest{UserID: input.UserID, PaginateParam: p}
	rsp, err := h.svc.ListUserLikeComments(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListUserLikeCommentsResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleListUserLikeTags(ctx context.Context, input *listWithPageInput) (*humadto.Output[protocol.ListUserLikeTagsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListUserLikeTagsRequest{UserID: input.UserID, PaginateParam: p}
	rsp, err := h.svc.ListUserLikeTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListUserLikeTagsResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleListImages(ctx context.Context, a *authHeader) (*humadto.Output[protocol.ListImagesResponse], error) {
	req := &protocol.ListImagesRequest{UserID: a.UserID}
	rsp, err := h.svc.ListImages(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListImagesResponse]{Body: *rsp}, nil
}

// HandleGetImage 返回预签名 URL，后续在路由中做 302 跳转
func (h *AssetHandlers) HandleGetImage(ctx context.Context, input *objectWithImageParamInput) (*humadto.Output[protocol.GetImageResponse], error) {
	req := &protocol.GetImageRequest{UserID: input.UserID, ImageName: input.ObjectName, Quality: input.Quality}
	rsp, err := h.svc.GetImage(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetImageResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleDeleteImage(ctx context.Context, input *objectPathInput) (*humadto.Output[protocol.DeleteImageResponse], error) {
	req := &protocol.DeleteImageRequest{UserID: input.UserID, ImageName: input.ObjectName}
	rsp, err := h.svc.DeleteImage(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteImageResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleListUserViewArticles(ctx context.Context, input *listWithPageInput) (*humadto.Output[protocol.ListUserViewArticlesResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListUserViewArticlesRequest{UserID: input.UserID, PaginateParam: p}
	rsp, err := h.svc.ListUserViewArticles(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListUserViewArticlesResponse]{Body: *rsp}, nil
}

func (h *AssetHandlers) HandleDeleteUserView(ctx context.Context, input *struct {
	authHeader
	humadto.ViewPath
},
) (*humadto.Output[protocol.DeleteUserViewResponse], error) {
	req := &protocol.DeleteUserViewRequest{UserID: input.UserID, ViewID: input.ViewID}
	rsp, err := h.svc.DeleteUserView(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteUserViewResponse]{Body: *rsp}, nil
}
