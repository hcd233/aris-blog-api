package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// ArticleHandler 文章处理器
//
//	author centonhuang
//	update 2025-10-30
type ArticleHandler interface {
	HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleResponse], error)
	HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoResponse], error)
	HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleResponse], error)
	HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleStatusResponse], error)
	HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HumaHTTPResponse[*dto.DeleteArticleResponse], error)
	HandleListArticles(ctx context.Context, req *dto.ListArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListArticlesResponse], error)
	HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleInfoBySlugRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoBySlugResponse], error)
}

type articleHandler struct {
	svc service.ArticleService
}

// NewArticleHandler 创建文章处理器
//
//	return ArticleHandler
//	author centonhuang
//	update 2025-10-30
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		svc: service.NewArticleService(),
	}
}

func (h *articleHandler) HandleCreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (*protocol.HumaHTTPResponse[*dto.CreateArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.CreateArticleRequest{
		UserID:     userID,
		Title:      req.Body.Title,
		Slug:       req.Body.Slug,
		CategoryID: req.Body.CategoryID,
		Tags:       req.Body.Tags,
	}

	svcRsp, err := h.svc.CreateArticle(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.CreateArticleResponse](nil, err)
	}

	rsp := &dto.CreateArticleResponse{
		Article: convertArticleToDTO(svcRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleHandler) HandleGetArticleInfo(ctx context.Context, req *dto.GetArticleInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetArticleInfoRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	svcRsp, err := h.svc.GetArticleInfo(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetArticleInfoResponse](nil, err)
	}

	rsp := &dto.GetArticleInfoResponse{
		Article: convertArticleToDTO(svcRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleHandler) HandleGetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleInfoBySlugRequest) (*protocol.HumaHTTPResponse[*dto.GetArticleInfoBySlugResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.GetArticleInfoBySlugRequest{
		UserID:      userID,
		AuthorName:  req.AuthorName,
		ArticleSlug: req.ArticleSlug,
	}

	svcRsp, err := h.svc.GetArticleInfoBySlug(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.GetArticleInfoBySlugResponse](nil, err)
	}

	rsp := &dto.GetArticleInfoBySlugResponse{
		Article: convertArticleToDTO(svcRsp.Article),
	}

	return util.WrapHTTPResponse(rsp, nil)
}

func (h *articleHandler) HandleUpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.UpdateArticleRequest{
		UserID:            userID,
		ArticleID:         req.ArticleID,
		UpdatedTitle:      req.Body.Title,
		UpdatedSlug:       req.Body.Slug,
		UpdatedCategoryID: req.Body.CategoryID,
	}

	_, err := h.svc.UpdateArticle(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.UpdateArticleResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.UpdateArticleResponse{}, nil)
}

func (h *articleHandler) HandleUpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (*protocol.HumaHTTPResponse[*dto.UpdateArticleStatusResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.UpdateArticleStatusRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}
	
	// 需要从string转换为ArticleStatus类型
	switch req.Body.Status {
	case "draft":
		svcReq.Status = "draft"
	case "publish":
		svcReq.Status = "publish"
	}

	_, err := h.svc.UpdateArticleStatus(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.UpdateArticleStatusResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.UpdateArticleStatusResponse{}, nil)
}

func (h *articleHandler) HandleDeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (*protocol.HumaHTTPResponse[*dto.DeleteArticleResponse], error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	svcReq := &protocol.DeleteArticleRequest{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	_, err := h.svc.DeleteArticle(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.DeleteArticleResponse](nil, err)
	}

	return util.WrapHTTPResponse(&dto.DeleteArticleResponse{}, nil)
}

func (h *articleHandler) HandleListArticles(ctx context.Context, req *dto.ListArticlesRequest) (*protocol.HumaHTTPResponse[*dto.ListArticlesResponse], error) {
	page := 1
	pageSize := 10
	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	svcReq := &protocol.ListArticlesRequest{
		PaginateParam: &protocol.PaginateParam{
			PageParam: &protocol.PageParam{
				Page:     page,
				PageSize: pageSize,
			},
		},
	}

	svcRsp, err := h.svc.ListArticles(ctx, svcReq)
	if err != nil {
		return util.WrapHTTPResponse[*dto.ListArticlesResponse](nil, err)
	}

	articles := make([]*dto.Article, len(svcRsp.Articles))
	for i, article := range svcRsp.Articles {
		articles[i] = convertArticleToDTO(article)
	}

	rsp := &dto.ListArticlesResponse{
		Articles: articles,
		PageInfo: &dto.PageInfo{
			Page:     svcRsp.PageInfo.Page,
			PageSize: svcRsp.PageInfo.PageSize,
			Total:    svcRsp.PageInfo.Total,
		},
	}

	return util.WrapHTTPResponse(rsp, nil)
}

// 辅助函数：转换Article类型
func convertArticleToDTO(article *protocol.Article) *dto.Article {
	if article == nil {
		return nil
	}

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
	for i, tag := range article.Tags {
		tags[i] = &dto.Tag{
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

	return &dto.Article{
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
