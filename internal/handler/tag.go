package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TagHandler 标签处理器
//
//	author centonhuang
//	update 2025-10-31 04:35:00
type TagHandler interface {
	HandleCreateTag(ctx context.Context, req *TagCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateTagResponse], error)
	HandleGetTagInfo(ctx context.Context, req *TagGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetTagInfoResponse], error)
	HandleUpdateTag(ctx context.Context, req *TagUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateTagResponse], error)
	HandleDeleteTag(ctx context.Context, req *TagDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteTagResponse], error)
	HandleListTags(ctx context.Context, req *TagListRequest) (*protocol.HumaHTTPResponse[*protocol.ListTagsResponse], error)
}

type tagHandler struct {
	svc service.TagService
}

// NewTagHandler 创建标签处理器
//
//	return TagHandler
//	author centonhuang
//	update 2025-10-31 04:35:00
func NewTagHandler() TagHandler {
	return &tagHandler{
		svc: service.NewTagService(),
	}
}

// TagPathParam 标签路径参数
type TagPathParam struct {
	TagID uint `path:"tagID" doc:"标签 ID"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	Body *protocol.CreateTagBody `json:"body" doc:"创建标签所需字段"`
}

// TagGetRequest 获取标签详情请求
type TagGetRequest struct {
	TagPathParam
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	TagPathParam
	Body *protocol.UpdateTagBody `json:"body" doc:"更新标签所需字段"`
}

// TagDeleteRequest 删除标签请求
type TagDeleteRequest struct {
	TagPathParam
}

// TagListRequest 标签列表请求
type TagListRequest struct {
	PaginationQuery
}

// HandleCreateTag 创建标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *TagCreateRequest
//	return *protocol.HumaHTTPResponse[*protocol.CreateTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:35:00
func (h *tagHandler) HandleCreateTag(ctx context.Context, req *TagCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateTagResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.CreateTagRequest{
		UserID:      userID,
		Name:        req.Body.Name,
		Slug:        req.Body.Slug,
		Description: req.Body.Description,
	}

	return util.WrapHTTPResponse(h.svc.CreateTag(ctx, serviceReq))
}

// HandleGetTagInfo 获取标签详情
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *TagGetRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetTagInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:35:00
func (h *tagHandler) HandleGetTagInfo(ctx context.Context, req *TagGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetTagInfoResponse], error) {
	serviceReq := &protocol.GetTagInfoRequest{
		TagID: req.TagID,
	}

	return util.WrapHTTPResponse(h.svc.GetTagInfo(ctx, serviceReq))
}

// HandleUpdateTag 更新标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *TagUpdateRequest
//	return *protocol.HumaHTTPResponse[*protocol.UpdateTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:35:00
func (h *tagHandler) HandleUpdateTag(ctx context.Context, req *TagUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateTagResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.UpdateTagRequest{
		UserID:      userID,
		TagID:       req.TagID,
		Name:        req.Body.Name,
		Slug:        req.Body.Slug,
		Description: req.Body.Description,
	}

	return util.WrapHTTPResponse(h.svc.UpdateTag(ctx, serviceReq))
}

// HandleDeleteTag 删除标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *TagDeleteRequest
//	return *protocol.HumaHTTPResponse[*protocol.DeleteTagResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:35:00
func (h *tagHandler) HandleDeleteTag(ctx context.Context, req *TagDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteTagResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.DeleteTagRequest{
		UserID: userID,
		TagID:  req.TagID,
	}

	return util.WrapHTTPResponse(h.svc.DeleteTag(ctx, serviceReq))
}

// HandleListTags 列出标签
//
//	receiver h *tagHandler
//	param ctx context.Context
//	param req *TagListRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListTagsResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:35:00
func (h *tagHandler) HandleListTags(ctx context.Context, req *TagListRequest) (*protocol.HumaHTTPResponse[*protocol.ListTagsResponse], error) {
	serviceReq := &protocol.ListTagsRequest{
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListTags(ctx, serviceReq))
}
