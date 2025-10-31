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
	HandleCreateTag(ctx context.Context, req *dto.CreateTagRequest) (*protocol.HumaHTTPResponse[*dto.CreateTagResponse], error)
	HandleGetTagInfo(ctx context.Context, req *dto.GetTagRequest) (*protocol.HumaHTTPResponse[*dto.GetTagResponse], error)
	HandleUpdateTag(ctx context.Context, req *dto.UpdateTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleDeleteTag(ctx context.Context, req *dto.DeleteTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
	HandleListTags(ctx context.Context, req *dto.ListTagRequest) (*protocol.HumaHTTPResponse[*dto.ListTagResponse], error)
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

func (h *tagHandler) HandleCreateTag(ctx context.Context, req *dto.CreateTagRequest) (*protocol.HumaHTTPResponse[*dto.CreateTagResponse], error) {
	return util.WrapHTTPResponse(h.svc.CreateTag(ctx, req))
}

func (h *tagHandler) HandleGetTagInfo(ctx context.Context, req *dto.GetTagRequest) (*protocol.HumaHTTPResponse[*dto.GetTagResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetTagInfo(ctx, req))
}

func (h *tagHandler) HandleUpdateTag(ctx context.Context, req *dto.UpdateTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateTag(ctx, req))
}

func (h *tagHandler) HandleDeleteTag(ctx context.Context, req *dto.DeleteTagRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.DeleteTag(ctx, req))
}

func (h *tagHandler) HandleListTags(ctx context.Context, req *dto.ListTagRequest) (*protocol.HumaHTTPResponse[*dto.ListTagResponse], error) {
	return util.WrapHTTPResponse(h.svc.ListTags(ctx, req))
}
