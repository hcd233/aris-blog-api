package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CategoryHandler 分类服务
type CategoryHandler interface {
	HandleCreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (*protocol.HumaHTTPResponse[*dto.CategoryCreateResponse], error)
	HandleGetCategoryInfo(ctx context.Context, req *dto.CategoryGetRequest) (*protocol.HumaHTTPResponse[*dto.CategoryGetResponse], error)
	HandleUpdateCategoryInfo(ctx context.Context, req *dto.CategoryUpdateRequest) (*protocol.HumaHTTPResponse[*dto.CategoryUpdateResponse], error)
	HandleDeleteCategory(ctx context.Context, req *dto.CategoryDeleteRequest) (*protocol.HumaHTTPResponse[*dto.CategoryDeleteResponse], error)
	HandleGetRootCategories(ctx context.Context, req *dto.CategoryGetRootRequest) (*protocol.HumaHTTPResponse[*dto.CategoryGetRootResponse], error)
	HandleListChildrenCategories(ctx context.Context, req *dto.CategoryListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*dto.CategoryListChildrenCategoriesResponse], error)
	HandleListChildrenArticles(ctx context.Context, req *dto.CategoryListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*dto.CategoryListChildrenArticlesResponse], error)
}

type categoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler() CategoryHandler {
	return &categoryHandler{
		svc: service.NewCategoryService(),
	}
}

func (h *categoryHandler) HandleCreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (*protocol.HumaHTTPResponse[*dto.CategoryCreateResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.CreateCategory(ctx, req))
}

func (h *categoryHandler) HandleGetCategoryInfo(ctx context.Context, req *dto.CategoryGetRequest) (*protocol.HumaHTTPResponse[*dto.CategoryGetResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.GetCategoryInfo(ctx, req))
}

func (h *categoryHandler) HandleUpdateCategoryInfo(ctx context.Context, req *dto.CategoryUpdateRequest) (*protocol.HumaHTTPResponse[*dto.CategoryUpdateResponse], error) {
	if req == nil {
		return nil, huma.Error400BadRequest("请求体不能为空")
	}

	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.UpdateCategory(ctx, req))
}

func (h *categoryHandler) HandleDeleteCategory(ctx context.Context, req *dto.CategoryDeleteRequest) (*protocol.HumaHTTPResponse[*dto.CategoryDeleteResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.DeleteCategory(ctx, req))
}

func (h *categoryHandler) HandleGetRootCategories(ctx context.Context, req *dto.CategoryGetRootRequest) (*protocol.HumaHTTPResponse[*dto.CategoryGetRootResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.GetRootCategory(ctx, req))
}

func (h *categoryHandler) HandleListChildrenCategories(ctx context.Context, req *dto.CategoryListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*dto.CategoryListChildrenCategoriesResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenCategories(ctx, req))
}

func (h *categoryHandler) HandleListChildrenArticles(ctx context.Context, req *dto.CategoryListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*dto.CategoryListChildrenArticlesResponse], error) {
	if userID, ok := UserIDFromCtx(ctx); ok {
		req.UserID = userID
	} else {
		return nil, huma.Error401Unauthorized("未登录或令牌无效")
	}

	return util.WrapHTTPResponse(h.svc.ListChildrenArticles(ctx, req))
}
