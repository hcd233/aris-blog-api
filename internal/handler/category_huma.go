package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// CreateCategoryHuma 创建分类（Huma 版本）
func CreateCategoryHuma(ctx context.Context, input *protocol.CreateCategoryInput) (*protocol.CategoryOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.CreateCategoryRequest{
		UserID:   userID,
		Name:     input.Body.Name,
		ParentID: input.Body.ParentID,
	}

	rsp, err := service.NewCategoryService().CreateCategory(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CategoryOutput{
		Body: *rsp.Category,
	}, nil
}

// GetCategoryInfoHuma 获取分类信息（Huma 版本）
func GetCategoryInfoHuma(ctx context.Context, input *protocol.CategoryInput) (*protocol.CategoryOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.GetCategoryInfoRequest{
		UserID:     userID,
		CategoryID: input.CategoryID,
	}

	rsp, err := service.NewCategoryService().GetCategoryInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CategoryOutput{
		Body: *rsp.Category,
	}, nil
}

// GetRootCategoriesHuma 获取根分类（Huma 版本）
func GetRootCategoriesHuma(ctx context.Context, input *struct{}) (*protocol.CategoryOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.GetRootCategoryRequest{
		UserID: userID,
	}

	rsp, err := service.NewCategoryService().GetRootCategory(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CategoryOutput{
		Body: *rsp.Category,
	}, nil
}

// UpdateCategoryHuma 更新分类（Huma 版本）
func UpdateCategoryHuma(ctx context.Context, input *protocol.UpdateCategoryInput) (*protocol.CategoryOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.UpdateCategoryRequest{
		UserID:     userID,
		CategoryID: input.CategoryID,
		Name:       input.Body.Name,
		ParentID:   input.Body.ParentID,
	}

	rsp, err := service.NewCategoryService().UpdateCategory(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CategoryOutput{
		Body: *rsp.Category,
	}, nil
}

// DeleteCategoryHuma 删除分类（Huma 版本）
func DeleteCategoryHuma(ctx context.Context, input *protocol.CategoryInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.DeleteCategoryRequest{
		UserID:     userID,
		CategoryID: input.CategoryID,
	}

	_, err := service.NewCategoryService().DeleteCategory(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}

// ListChildrenCategoriesHuma 列出子分类（Huma 版本）
func ListChildrenCategoriesHuma(ctx context.Context, input *protocol.CategoryInput) (*protocol.CategoryListOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.ListChildrenCategoriesRequest{
		UserID:     userID,
		CategoryID: input.CategoryID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     1,
				PageSize: 50,
			},
		},
	}

	rsp, err := service.NewCategoryService().ListChildrenCategories(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	categories := make([]protocol.Category, len(rsp.Categories))
	for i, c := range rsp.Categories {
		categories[i] = *c
	}
	return &protocol.CategoryListOutput{
		Body: struct {
			Categories []protocol.Category `json:"categories"`
			PageInfo   protocol.PageInfo   `json:"pageInfo"`
		}{
			Categories: categories,
			PageInfo:   *rsp.PageInfo,
		},
	}, nil
}
