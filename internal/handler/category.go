package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CategoryHandler 分类处理器
//
//	author centonhuang
//	update 2025-10-30
type CategoryHandler interface {
	HandleCreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*protocol.HumaHTTPResponse[*dto.CreateCategoryResponse], error)
	HandleGetCategoryInfo(ctx context.Context, req *dto.GetCategoryInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetCategoryInfoResponse], error)
	HandleUpdateCategoryInfo(ctx context.Context, req *dto.UpdateCategoryRequest) (*protocol.HumaHTTPResponse[*dto.UpdateCategoryResponse], error)
	HandleDeleteCategory(ctx context.Context, req *dto.DeleteCategoryRequest) (*protocol.HumaHTTPResponse[*dto.DeleteCategoryResponse], error)
	HandleGetRootCategories(ctx context.Context, req *dto.GetRootCategoryRequest) (*protocol.HumaHTTPResponse[*dto.GetRootCategoryResponse], error)
	HandleListChildrenCategories(ctx context.Context, req *dto.ListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenCategoriesResponse], error)
	HandleListChildrenArticles(ctx context.Context, req *dto.ListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenArticlesResponse], error)
}

type categoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 创建分类处理器
//
//	return CategoryHandler
//	author centonhuang
//	update 2025-10-30
func NewCategoryHandler() CategoryHandler {
	return &categoryHandler{
		svc: service.NewCategoryService(),
	}
}

func (h *categoryHandler) HandleCreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*protocol.HumaHTTPResponse[*dto.CreateCategoryResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.CreateCategoryRequest{
		UserID:   userID,
		Name:     req.Body.Name,
		ParentID: req.Body.ParentID,
	}

	svcRsp, err := h.svc.CreateCategory(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.CreateCategoryResponse](nil, err)
	}

	rsp := &dto.CreateCategoryResponse{
		Category: &dto.Category{
			CategoryID: svcRsp.Category.CategoryID,
			Name:       svcRsp.Category.Name,
			ParentID:   svcRsp.Category.ParentID,
			CreatedAt:  svcRsp.Category.CreatedAt,
			UpdatedAt:  svcRsp.Category.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *categoryHandler) HandleGetCategoryInfo(ctx context.Context, req *dto.GetCategoryInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetCategoryInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetCategoryInfoRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
	}

	svcRsp, err := h.svc.GetCategoryInfo(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetCategoryInfoResponse](nil, err)
	}

	rsp := &dto.GetCategoryInfoResponse{
		Category: &dto.Category{
			CategoryID: svcRsp.Category.CategoryID,
			Name:       svcRsp.Category.Name,
			ParentID:   svcRsp.Category.ParentID,
			CreatedAt:  svcRsp.Category.CreatedAt,
			UpdatedAt:  svcRsp.Category.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *categoryHandler) HandleUpdateCategoryInfo(ctx context.Context, req *dto.UpdateCategoryRequest) (*protocol.HumaHTTPResponse[*dto.UpdateCategoryResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.UpdateCategoryRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Name:       req.Body.Name,
		ParentID:   req.Body.ParentID,
	}

	svcRsp, err := h.svc.UpdateCategory(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.UpdateCategoryResponse](nil, err)
	}

	rsp := &dto.UpdateCategoryResponse{
		Category: &dto.Category{
			CategoryID: svcRsp.Category.CategoryID,
			Name:       svcRsp.Category.Name,
			ParentID:   svcRsp.Category.ParentID,
			CreatedAt:  svcRsp.Category.CreatedAt,
			UpdatedAt:  svcRsp.Category.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *categoryHandler) HandleDeleteCategory(ctx context.Context, req *dto.DeleteCategoryRequest) (*protocol.HumaHTTPResponse[*dto.DeleteCategoryResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.DeleteCategoryRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
	}

	_, err := h.svc.DeleteCategory(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.DeleteCategoryResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.DeleteCategoryResponse{}, nil)
}

func (h *categoryHandler) HandleGetRootCategories(ctx context.Context, req *dto.GetRootCategoryRequest) (*protocol.HumaHTTPResponse[*dto.GetRootCategoryResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetRootCategoryRequest{
		UserID: userID,
	}

	svcRsp, err := h.svc.GetRootCategory(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetRootCategoryResponse](nil, err)
	}

	rsp := &dto.GetRootCategoryResponse{
		Category: &dto.Category{
			CategoryID: svcRsp.Category.CategoryID,
			Name:       svcRsp.Category.Name,
			ParentID:   svcRsp.Category.ParentID,
			CreatedAt:  svcRsp.Category.CreatedAt,
			UpdatedAt:  svcRsp.Category.UpdatedAt,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *categoryHandler) HandleListChildrenCategories(ctx context.Context, req *dto.ListChildrenCategoriesRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenCategoriesResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListChildrenCategoriesRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListChildrenCategories(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListChildrenCategoriesResponse](nil, err)
	}

	categories := make([]*dto.Category, len(svcRsp.Categories))
	for i, cat := range svcRsp.Categories {
		categories[i] = &dto.Category{
			CategoryID: cat.CategoryID,
			Name:       cat.Name,
			ParentID:   cat.ParentID,
			CreatedAt:  cat.CreatedAt,
			UpdatedAt:  cat.UpdatedAt,
		}
	}

	rsp := &dto.ListChildrenCategoriesResponse{
		Categories: categories,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *categoryHandler) HandleListChildrenArticles(ctx context.Context, req *dto.ListChildrenArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListChildrenArticlesResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListChildrenArticlesRequest{
		UserID:     userID,
		CategoryID: req.CategoryID,
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListChildrenArticles(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListChildrenArticlesResponse](nil, err)
	}

	articles := make([]*dto.Article, len(svcRsp.Articles))
	for i, article := range svcRsp.Articles {
		var user *dto.User
		if article.User != nil {
			user = &dto.User{
				UserID:    article.User.UserID,
				Name:      article.User.Name,
				Email:     article.User.Email,
				Avatar:    article.User.Avatar,
				CreatedAt: article.User.CreatedAt,
				LastLogin: article.User.LastLogin,
			}
		}

		var category *dto.Category
		if article.Category != nil {
			category = &dto.Category{
				CategoryID: article.Category.CategoryID,
				Name:       article.Category.Name,
				ParentID:   article.Category.ParentID,
				CreatedAt:  article.Category.CreatedAt,
				UpdatedAt:  article.Category.UpdatedAt,
			}
		}

		tags := make([]*dto.Tag, len(article.Tags))
		for j, tag := range article.Tags {
			tags[j] = &dto.Tag{
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

		articles[i] = &dto.Article{
			ArticleID:   article.ArticleID,
			Title:       article.Title,
			Slug:        article.Slug,
			Status:      article.Status,
			User:        user,
			Category:    category,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
			PublishedAt: article.PublishedAt,
			Likes:       article.Likes,
			Views:       article.Views,
			Tags:        tags,
			Comments:    article.Comments,
		}
	}

	rsp := &dto.ListChildrenArticlesResponse{
		Articles: articles,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}
