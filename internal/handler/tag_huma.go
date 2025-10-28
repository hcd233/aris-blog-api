package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// CreateTagHuma 创建标签（Huma 版本）
func CreateTagHuma(ctx context.Context, input *protocol.CreateTagInput) (*protocol.TagOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.CreateTagRequest{
		UserID:      userID,
		Name:        input.Body.Name,
		Slug:        input.Body.Slug,
		Description: input.Body.Description,
	}

	rsp, err := service.NewTagService().CreateTag(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.TagOutput{
		Body: *rsp.Tag,
	}, nil
}

// GetTagInfoHuma 获取标签信息（Huma 版本）
func GetTagInfoHuma(ctx context.Context, input *protocol.TagInput) (*protocol.TagOutput, error) {
	req := &protocol.GetTagInfoRequest{
		TagID: input.TagID,
	}

	rsp, err := service.NewTagService().GetTagInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.TagOutput{
		Body: *rsp.Tag,
	}, nil
}

// UpdateTagHuma 更新标签（Huma 版本）
func UpdateTagHuma(ctx context.Context, input *protocol.UpdateTagInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.UpdateTagRequest{
		UserID:      userID,
		TagID:       input.TagID,
		Name:        input.Body.Name,
		Slug:        input.Body.Slug,
		Description: input.Body.Description,
	}

	_, err := service.NewTagService().UpdateTag(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// DeleteTagHuma 删除标签（Huma 版本）
func DeleteTagHuma(ctx context.Context, input *protocol.TagInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.DeleteTagRequest{
		UserID: userID,
		TagID:  input.TagID,
	}

	_, err := service.NewTagService().DeleteTag(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// ListTagsHuma 列出标签（Huma 版本）
func ListTagsHuma(ctx context.Context, input *protocol.PaginatedSearchParams) (*protocol.TagListOutput, error) {
	req := &protocol.ListTagsRequest{
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     input.Page,
				PageSize: input.PageSize,
			},
			QueryParam: &protocol.QueryParam{
				Query: input.Query,
			},
		},
	}

	rsp, err := service.NewTagService().ListTags(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	tags := make([]protocol.Tag, len(rsp.Tags))
	for i, t := range rsp.Tags {
		tags[i] = *t
	}
	return &protocol.TagListOutput{
		Body: struct {
			Tags     []protocol.Tag    `json:"tags"`
			PageInfo protocol.PageInfo `json:"pageInfo"`
		}{
			Tags:     tags,
			PageInfo: *rsp.PageInfo,
		},
	}, nil
}
