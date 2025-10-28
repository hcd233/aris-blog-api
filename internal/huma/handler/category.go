package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type CategoryHandlers struct {
	svc service.CategoryService
}

func NewCategoryHandlers() *CategoryHandlers {
	return &CategoryHandlers{svc: service.NewCategoryService()}
}

type (
	createCategoryInput struct {
		authHeader
		humadto.CreateCategoryInput
	}

	categoryPathInput struct {
		authHeader
		humadto.CategoryPath
	}

	updateCategoryInput struct {
		authHeader
		humadto.CategoryPath
		humadto.UpdateCategoryInput
	}

	listChildrenInput struct {
		authHeader
		humadto.CategoryPath
		humadto.PaginateParam
	}
)

func (h *CategoryHandlers) HandleCreateCategory(ctx context.Context, input *createCategoryInput) (*humadto.Output[protocol.CreateCategoryResponse], error) {
	req := &protocol.CreateCategoryRequest{
		UserID:   input.UserID,
		Name:     input.Body.Name,
		ParentID: input.Body.ParentID,
	}
	rsp, err := h.svc.CreateCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CreateCategoryResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleGetCategoryInfo(ctx context.Context, input *categoryPathInput) (*humadto.Output[protocol.GetCategoryInfoResponse], error) {
	req := &protocol.GetCategoryInfoRequest{UserID: input.UserID, CategoryID: input.CategoryID}
	rsp, err := h.svc.GetCategoryInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetCategoryInfoResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleGetRootCategories(ctx context.Context, _ *authHeader) (*humadto.Output[protocol.GetRootCategoryResponse], error) {
	req := &protocol.GetRootCategoryRequest{UserID: 0}
	rsp, err := h.svc.GetRootCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetRootCategoryResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleUpdateCategoryInfo(ctx context.Context, input *updateCategoryInput) (*humadto.Output[protocol.UpdateCategoryResponse], error) {
	req := &protocol.UpdateCategoryRequest{
		UserID:     input.UserID,
		CategoryID: input.CategoryID,
		Name:       input.Body.Name,
		ParentID:   input.Body.ParentID,
	}
	rsp, err := h.svc.UpdateCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.UpdateCategoryResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleDeleteCategory(ctx context.Context, input *categoryPathInput) (*humadto.Output[protocol.DeleteCategoryResponse], error) {
	req := &protocol.DeleteCategoryRequest{UserID: input.UserID, CategoryID: input.CategoryID}
	rsp, err := h.svc.DeleteCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.DeleteCategoryResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleListChildrenCategories(ctx context.Context, input *listChildrenInput) (*humadto.Output[protocol.ListChildrenCategoriesResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListChildrenCategoriesRequest{UserID: input.UserID, CategoryID: input.CategoryID, PaginateParam: p}
	rsp, err := h.svc.ListChildrenCategories(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListChildrenCategoriesResponse]{Body: *rsp}, nil
}

func (h *CategoryHandlers) HandleListChildrenArticles(ctx context.Context, input *listChildrenInput) (*humadto.Output[protocol.ListChildrenArticlesResponse], error) {
	p := &protocol.PaginateParam{PageParam: &protocol.PageParam{}, QueryParam: &protocol.QueryParam{}}
	if input.PageParam != nil {
		p.PageParam = &protocol.PageParam{Page: input.Page, PageSize: input.PageSize}
	}
	if input.QueryParam != nil {
		p.QueryParam = &protocol.QueryParam{Query: input.Query}
	}
	req := &protocol.ListChildrenArticlesRequest{UserID: input.UserID, CategoryID: input.CategoryID, PaginateParam: p}
	rsp, err := h.svc.ListChildrenArticles(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.ListChildrenArticlesResponse]{Body: *rsp}, nil
}
