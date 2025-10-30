package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CategoryHandler 分类服务
//
//	author centonhuang
//	update 2025-10-31 04:20:00
type CategoryHandler interface {
	HandleCreateCategory(ctx context.Context, req *CategoryCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateCategoryResponse], error)
	HandleGetCategoryInfo(ctx context.Context, req *CategoryGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetCategoryInfoResponse], error)
	HandleUpdateCategoryInfo(ctx context.Context, req *CategoryUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateCategoryResponse], error)
	HandleDeleteCategory(ctx context.Context, req *CategoryDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteCategoryResponse], error)
	HandleGetRootCategories(ctx context.Context, req *CategoryGetRootRequest) (*protocol.HumaHTTPResponse[*protocol.GetRootCategoryResponse], error)
	HandleListChildrenCategories(ctx context.Context, req *CategoryListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenCategoriesResponse], error)
	HandleListChildrenArticles(ctx context.Context, req *CategoryListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenArticlesResponse], error)
}

type categoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 创建分类处理器
//
//	return CategoryHandler
//	author centonhuang
//	update 2025-10-31 04:20:00
func NewCategoryHandler() CategoryHandler {
	return &categoryHandler{
		svc: service.NewCategoryService(),
	}
}

// CategoryPathParam 分类路径参数
type CategoryPathParam struct {
	CategoryID uint `path:"categoryID" doc:"分类 ID"`
}

// CategoryCreateRequest 创建分类请求
type CategoryCreateRequest struct {
	Body *protocol.CreateCategoryBody `json:"body" doc:"创建分类所需字段"`
}

// CategoryGetRequest 获取分类详情请求
type CategoryGetRequest struct {
	CategoryPathParam
}

// CategoryUpdateRequest 更新分类请求
type CategoryUpdateRequest struct {
	CategoryPathParam
	Body *protocol.UpdateCategoryBody `json:"body" doc:"更新分类所需字段"`
}

// CategoryDeleteRequest 删除分类请求
type CategoryDeleteRequest struct {
	CategoryPathParam
}

// CategoryGetRootRequest 获取根分类请求
type CategoryGetRootRequest struct{}

// CategoryListChildrenCategoriesRequest 列出子分类请求
type CategoryListChildrenCategoriesRequest struct {
	CategoryPathParam
	PaginationQuery
}

// CategoryListChildrenArticlesRequest 列出子文章请求
type CategoryListChildrenArticlesRequest struct {
	CategoryPathParam
	PaginationQuery
}

// HandleCreateCategory 创建分类
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryCreateRequest
//	return *protocol.HumaHTTPResponse[*protocol.CreateCategoryResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleCreateCategory(ctx context.Context, req *CategoryCreateRequest) (*protocol.HumaHTTPResponse[*protocol.CreateCategoryResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.CreateCategoryRequest{
		UserID:   userID,
		Name:     req.Body.Name,
		ParentID: req.Body.ParentID,
	}

	return util.WrapHTTPResponse(h.svc.CreateCategory(ctx, serviceReq))
}

// HandleGetCategoryInfo 获取分类详情
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryGetRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetCategoryInfoResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleGetCategoryInfo(ctx context.Context, req *CategoryGetRequest) (*protocol.HumaHTTPResponse[*protocol.GetCategoryInfoResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetCategoryInfoRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
	}

	return util.WrapHTTPResponse(h.svc.GetCategoryInfo(ctx, serviceReq))
}

// HandleGetRootCategories 获取根分类
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param _ *CategoryGetRootRequest
//	return *protocol.HumaHTTPResponse[*protocol.GetRootCategoryResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleGetRootCategories(ctx context.Context, _ *CategoryGetRootRequest) (*protocol.HumaHTTPResponse[*protocol.GetRootCategoryResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.GetRootCategoryRequest{
		UserID: userID,
	}

	return util.WrapHTTPResponse(h.svc.GetRootCategory(ctx, serviceReq))
}

// HandleUpdateCategoryInfo 更新分类
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryUpdateRequest
//	return *protocol.HumaHTTPResponse[*protocol.UpdateCategoryResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleUpdateCategoryInfo(ctx context.Context, req *CategoryUpdateRequest) (*protocol.HumaHTTPResponse[*protocol.UpdateCategoryResponse], error) {
	if req == nil || req.Body == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.UpdateCategoryRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Name:       req.Body.Name,
		ParentID:   req.Body.ParentID,
	}

	return util.WrapHTTPResponse(h.svc.UpdateCategory(ctx, serviceReq))
}

// HandleDeleteCategory 删除分类
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryDeleteRequest
//	return *protocol.HumaHTTPResponse[*protocol.DeleteCategoryResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleDeleteCategory(ctx context.Context, req *CategoryDeleteRequest) (*protocol.HumaHTTPResponse[*protocol.DeleteCategoryResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.DeleteCategoryRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
	}

	return util.WrapHTTPResponse(h.svc.DeleteCategory(ctx, serviceReq))
}

// HandleListChildrenCategories 列出子分类
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryListChildrenCategoriesRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListChildrenCategoriesResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleListChildrenCategories(ctx context.Context, req *CategoryListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenCategoriesResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.ListChildrenCategoriesRequest{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenCategories(ctx, serviceReq))
}

// HandleListChildrenArticles 列出子文章
//
//	receiver h *categoryHandler
//	param ctx context.Context
//	param req *CategoryListChildrenArticlesRequest
//	return *protocol.HumaHTTPResponse[*protocol.ListChildrenArticlesResponse]
//	return error
//	author centonhuang
//	update 2025-10-31 04:20:00
func (h *categoryHandler) HandleListChildrenArticles(ctx context.Context, req *CategoryListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*protocol.ListChildrenArticlesResponse], error) {
	userID, ok := UserIDFromCtx(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	serviceReq := &protocol.ListChildrenArticlesRequest{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		PaginateParam: req.PaginationQuery.ToPaginateParam(),
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenArticles(ctx, serviceReq))
}
