package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type TagHandlers struct{ svc service.TagService }

func NewTagHandlers() *TagHandlers { return &TagHandlers{svc: service.NewTagService()} }

type (
	createTagInput struct {
		authHeader
		humadto.CreateTagInput
	}
	tagPathInput struct {
		authHeader
		humadto.TagPath
	}
	updateTagInput struct {
		authHeader
		humadto.TagPath
		humadto.UpdateTagInput
	}
	listTagsInput struct{ humadto.PaginateParam }
)

func (h *TagHandlers) HandleCreateTag(ctx context.Context, input *createTagInput) (*humadto.Output[protocol.CreateTagResponse], error) {
	req := protocol.CreateTagRequest{UserID: input.UserID, Name: input.Body.Name, Slug: input.Body.Slug, Description: input.Body.Description}
	rsp, err := h.svc.CreateTag(ctx, &req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreateTagResponse]{Body: *rsp}, nil
}

func (h *TagHandlers) HandleGetTagInfo(ctx context.Context, input *tagPathInput) (*humadto.Output[protocol.GetTagInfoResponse], error) {
	req := &protocol.GetTagInfoRequest{TagID: input.TagID}
	rsp, err := h.svc.GetTagInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetTagInfoResponse]{Body: *rsp}, nil
}

func (h *TagHandlers) HandleUpdateTag(ctx context.Context, input *updateTagInput) (*humadto.Output[protocol.UpdateTagResponse], error) {
	req := &protocol.UpdateTagRequest{UserID: input.UserID, TagID: input.TagID, Name: input.Body.Name, Slug: input.Body.Slug, Description: input.Body.Description}
	rsp, err := h.svc.UpdateTag(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.UpdateTagResponse]{Body: *rsp}, nil
}

func (h *TagHandlers) HandleDeleteTag(ctx context.Context, input *tagPathInput) (*humadto.Output[protocol.DeleteTagResponse], error) {
	req := &protocol.DeleteTagRequest{UserID: input.UserID, TagID: input.TagID}
	rsp, err := h.svc.DeleteTag(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteTagResponse]{Body: *rsp}, nil
}

func (h *TagHandlers) HandleListTags(ctx context.Context, input *listTagsInput) (*humadto.Output[protocol.ListTagsResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListTagsRequest{PaginateParam: p}
	rsp, err := h.svc.ListTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListTagsResponse]{Body: *rsp}, nil
}
