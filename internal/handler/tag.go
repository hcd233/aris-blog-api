package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TagHandler 标签处理器
type TagHandler interface {
	HandleCreateTag(ctx context.Context, req *dto.TagCreateRequest) (*protocol.HumaHTTPResponse[*dto.TagCreateResponse], error)
	HandleGetTagInfo(ctx context.Context, req *dto.TagGetRequest) (*protocol.HumaHTTPResponse[*dto.TagGetResponse], error)
	HandleUpdateTag(ctx context.Context, req *dto.TagUpdateRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleDeleteTag(ctx context.Context, req *dto.TagDeleteRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleListTags(ctx context.Context, req *dto.TagListRequest) (*protocol.HumaHTTPResponse[*dto.TagListResponse], error)
}

type tagHandler struct {
	svc service.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler() TagHandler {
	return &tagHandler{
		svc: service.NewTagService(),
	}
}

func (h *tagHandler) HandleCreateTag(ctx context.Context, req *dto.TagCreateRequest) (*protocol.HumaHTTPResponse[*dto.TagCreateResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateTag(ctx, req))
}

func (h *tagHandler) HandleGetTagInfo(ctx context.Context, req *dto.TagGetRequest) (*protocol.HumaHTTPResponse[*dto.TagGetResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetTagInfo(ctx, req))
}

func (h *tagHandler) HandleUpdateTag(ctx context.Context, req *dto.TagUpdateRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateTag(ctx, req))
}

func (h *tagHandler) HandleDeleteTag(ctx context.Context, req *dto.TagDeleteRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteTag(ctx, req))
}

func (h *tagHandler) HandleListTags(ctx context.Context, req *dto.TagListRequest) (*protocol.HumaHTTPResponse[*dto.TagListResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListTags(ctx, req))
}
