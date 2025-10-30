package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TagHandler 标签处理器
//
//	author centonhuang
//	update 2025-10-30
type TagHandler interface {
	HandleCreateTag(ctx context.Context, req *dto.CreateTagRequest) (*protocol.HumaHTTPResponse[*dto.CreateTagResponse], error)
	HandleGetTagInfo(ctx context.Context, req *dto.GetTagInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetTagInfoResponse], error)
	HandleUpdateTag(ctx context.Context, req *dto.UpdateTagRequest) (*protocol.HumaHTTPResponse[*dto.UpdateTagResponse], error)
	HandleDeleteTag(ctx context.Context, req *dto.DeleteTagRequest) (*protocol.HumaHTTPResponse[*dto.DeleteTagResponse], error)
	HandleListTags(ctx context.Context, req *dto.ListTagsRequest) (*protocol.HumaHTTPResponse[*dto.ListTagsResponse], error)
}

type tagHandler struct {
	svc service.TagService
}

// NewTagHandler 创建标签处理器
//
//	return TagHandler
//	author centonhuang
//	update 2025-10-30
func NewTagHandler() TagHandler {
	return &tagHandler{
		svc: service.NewTagService(),
	}
}

// HandleCreateTag 创建标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *dto.CreateTagRequest
//	return *protocol.HumaHTTPResponse[*dto.CreateTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-30
func (h *tagHandler) HandleCreateTag(ctx context.Context, req *dto.CreateTagRequest) (*protocol.HumaHTTPResponse[*dto.CreateTagResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.CreateTagRequest{
		UserID:      userID,
		Name:        req.Body.Name,
		Slug:        req.Body.Slug,
		Description: req.Body.Description,
	}

	svcRsp, err := h.svc.CreateTag(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.CreateTagResponse](nil, err)
	}

	rsp := &dto.CreateTagResponse{
		Tag: &dto.Tag{
			TagID:       svcRsp.Tag.TagID,
			Name:        svcRsp.Tag.Name,
			Slug:        svcRsp.Tag.Slug,
			Description: svcRsp.Tag.Description,
			UserID:      svcRsp.Tag.UserID,
			CreatedAt:   svcRsp.Tag.CreatedAt,
			UpdatedAt:   svcRsp.Tag.UpdatedAt,
			Likes:       svcRsp.Tag.Likes,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleGetTagInfo 获取标签信息
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *dto.GetTagInfoRequest
//	return *protocol.HumaHTTPResponse[*dto.GetTagInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-30
func (h *tagHandler) HandleGetTagInfo(ctx context.Context, req *dto.GetTagInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetTagInfoResponse], error) {
	svcReq := &protocol.GetTagInfoRequest{
		TagID: req.TagID,
	}

	svcRsp, err := h.svc.GetTagInfo(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetTagInfoResponse](nil, err)
	}

	rsp := &dto.GetTagInfoResponse{
		Tag: &dto.Tag{
			TagID:       svcRsp.Tag.TagID,
			Name:        svcRsp.Tag.Name,
			Slug:        svcRsp.Tag.Slug,
			Description: svcRsp.Tag.Description,
			UserID:      svcRsp.Tag.UserID,
			CreatedAt:   svcRsp.Tag.CreatedAt,
			UpdatedAt:   svcRsp.Tag.UpdatedAt,
			Likes:       svcRsp.Tag.Likes,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// HandleUpdateTag 更新标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *dto.UpdateTagRequest
//	return *protocol.HumaHTTPResponse[*dto.UpdateTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-30
func (h *tagHandler) HandleUpdateTag(ctx context.Context, req *dto.UpdateTagRequest) (*protocol.HumaHTTPResponse[*dto.UpdateTagResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.UpdateTagRequest{
		UserID:      userID,
		TagID:       req.TagID,
		Name:        req.Body.Name,
		Slug:        req.Body.Slug,
		Description: req.Body.Description,
	}

	_, err := h.svc.UpdateTag(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.UpdateTagResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.UpdateTagResponse{}, nil)
}

// HandleDeleteTag 删除标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *dto.DeleteTagRequest
//	return *protocol.HumaHTTPResponse[*dto.DeleteTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-30
func (h *tagHandler) HandleDeleteTag(ctx context.Context, req *dto.DeleteTagRequest) (*protocol.HumaHTTPResponse[*dto.DeleteTagResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.DeleteTagRequest{
		UserID: userID,
		TagID:  req.TagID,
	}

	_, err := h.svc.DeleteTag(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.DeleteTagResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.DeleteTagResponse{}, nil)
}

// HandleListTags 列出标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *dto.ListTagsRequest
//	return *protocol.HumaHTTPResponse[*dto.ListTagsResponse]
//	return error
//	author centonhuang
//	update 2025-10-30
func (h *tagHandler) HandleListTags(ctx context.Context, req *dto.ListTagsRequest) (*protocol.HumaHTTPResponse[*dto.ListTagsResponse], error) {
	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListTagsRequest{
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListTags(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListTagsResponse](nil, err)
	}

	tags := make([]*dto.Tag, len(svcRsp.Tags))
	for i, tag := range svcRsp.Tags {
		tags[i] = &dto.Tag{
			TagID:       tag.TagID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			UserID:      tag.UserID,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
			Likes:       tag.Likes,
		}
	}

	rsp := &dto.ListTagsResponse{
		Tags: tags,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}
