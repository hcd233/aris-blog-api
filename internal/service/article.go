package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleService 文章服务
type ArticleService interface {
	CreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (rsp *dto.CreateArticleResponse, err error)
	GetArticleInfo(ctx context.Context, req *dto.GetArticleRequest) (rsp *dto.GetArticleResponse, err error)
	GetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleBySlugRequest) (rsp *dto.GetArticleBySlugResponse, err error)
	UpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (rsp *dto.EmptyResponse, err error)
	UpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (rsp *dto.EmptyResponse, err error)
	DeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (rsp *dto.EmptyResponse, err error)
	ListArticles(ctx context.Context, req *dto.ListArticleRequest) (rsp *dto.ListArticleResponse, err error)
}

type articleService struct {
	userDAO           *dao.UserDAO
	tagDAO            *dao.TagDAO
	categoryDAO       *dao.CategoryDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
}

// NewArticleService 创建文章服务
func NewArticleService() ArticleService {
	return &articleService{
		userDAO:           dao.GetUserDAO(),
		tagDAO:            dao.GetTagDAO(),
		categoryDAO:       dao.GetCategoryDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
	}
}

// CreateArticle 创建文章
func (s *articleService) CreateArticle(ctx context.Context, req *dto.CreateArticleRequest) (rsp *dto.CreateArticleResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[UserService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.CreateArticleResponse{}

	db := database.GetDBInstance(ctx)

	tags := []model.Tag{}
	tagChan, errChan := make(chan *model.Tag, len(req.Body.Tags)), make(chan error, len(req.Body.Tags))

	var wg sync.WaitGroup
	wg.Add(len(req.Body.Tags))

	getTagFunc := func(tagSlug string) {
		defer wg.Done()
		tag, err := s.tagDAO.GetBySlug(db, tagSlug, []string{"id", "slug"}, []string{})
		if err != nil {
			errChan <- err
			return
		}
		tagChan <- tag
	}

	for _, tagSlug := range req.Body.Tags {
		go getTagFunc(tagSlug)
	}

	wg.Wait()
	close(tagChan)
	close(errChan)

	if len(errChan) > 0 {
		err := <-errChan
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] tag not found", zap.Error(err))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get tag", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	for tag := range tagChan {
		tags = append(tags, *tag)
	}

	category, err := s.categoryDAO.GetByID(db, req.Body.CategoryID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] category not found",
				zap.Uint("categoryID", req.Body.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get category", zap.Uint("categoryID", req.Body.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article := &model.Article{
		UserID:   userID,
		Status:   model.ArticleStatusDraft,
		Title:    req.Body.Title,
		Slug:     req.Body.Slug,
		Tags:     tags,
		Category: category,
		Comments: []model.Comment{},
		Versions: []model.ArticleVersion{},
	}

	if err := s.articleDAO.Create(db, article); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Error("[ArticleService] article slug duplicated",
				zap.String("title", article.Title),
				zap.String("slug", article.Slug))
			return nil, protocol.ErrDataExists
		}
		logger.Error("[ArticleService] failed to create article",
			zap.String("title", article.Title),
			zap.String("slug", article.Slug),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Article = &dto.Article{
		ArticleID:   article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      string(article.Status),
		User:        nil,
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags:        nil,
		Comments:    len(article.Comments),
	}

	return rsp, nil
}

// GetArticleInfo 获取文章信息
func (s *articleService) GetArticleInfo(ctx context.Context, req *dto.GetArticleRequest) (rsp *dto.GetArticleResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil {
		logger.Error("[ArticleService] request is nil")
		return nil, protocol.ErrBadRequest
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.GetArticleResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{
		"id", "slug", "title", "status", "user_id", "category_id",
		"created_at", "updated_at", "published_at",
		"likes", "views",
	}, []string{"User", "Category", "Tags", "Comments"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != userID && article.Status != model.ArticleStatusPublish {
		logger.Error("[ArticleService] no permission to get article",
			zap.Uint("articleID", req.ArticleID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Article = s.buildArticleDTO(article)

	return rsp, nil
}

// GetArticleInfoBySlug 通过别名获取文章信息
func (s *articleService) GetArticleInfoBySlug(ctx context.Context, req *dto.GetArticleBySlugRequest) (rsp *dto.GetArticleBySlugResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil {
		logger.Error("[ArticleService] request is nil")
		return nil, protocol.ErrBadRequest
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.GetArticleBySlugResponse{}

	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.GetByName(db, req.AuthorName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] user not found",
				zap.String("authorName", req.AuthorName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get user",
			zap.String("authorName", req.AuthorName),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(db, req.ArticleSlug, user.ID, []string{
		"id", "slug", "title", "status", "user_id", "category_id",
		"created_at", "updated_at", "published_at",
		"likes", "views",
	}, []string{"User", "Category", "Tags", "Comments"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] article not found",
				zap.String("slug", req.ArticleSlug))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get article",
			zap.String("slug", req.ArticleSlug),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != userID && article.Status != model.ArticleStatusPublish {
		logger.Error("[ArticleService] no permission to get article",
			zap.String("ArticleSlug", req.ArticleSlug),
			zap.String("AuthorName", req.AuthorName))
		return nil, protocol.ErrNoPermission
	}

	rsp.Article = s.buildArticleDTO(article)

	return rsp, nil
}

// UpdateArticle 更新文章
func (s *articleService) UpdateArticle(ctx context.Context, req *dto.UpdateArticleRequest) (rsp *dto.EmptyResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[ArticleService] request is nil")
		return nil, protocol.ErrBadRequest
	}

	rsp = &dto.EmptyResponse{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	db := database.GetDBInstance(ctx)

	updateFields := make(map[string]interface{})
	if req.Body.Title != "" {
		updateFields["title"] = req.Body.Title
	}
	if req.Body.Slug != "" {
		updateFields["slug"] = req.Body.Slug
	}
	if req.Body.CategoryID != 0 {
		updateFields["category_id"] = req.Body.CategoryID
	}

	if len(updateFields) == 0 {
		logger.Warn("[ArticleService] no fields to update",
			zap.Uint("articleID", req.ArticleID))
		return rsp, nil
	}

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != userID {
		logger.Error("[ArticleService] no permission to update article",
			zap.Uint("articleID", article.ID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.articleDAO.Update(db, article, updateFields); err != nil {
		logger.Error("[ArticleService] failed to update article",
			zap.Uint("articleID", article.ID),
			zap.Any("updateFields", updateFields),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// UpdateArticleStatus 更新文章状态
func (s *articleService) UpdateArticleStatus(ctx context.Context, req *dto.UpdateArticleStatusRequest) (rsp *dto.EmptyResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[ArticleService] request is nil")
		return nil, protocol.ErrBadRequest
	}
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.EmptyResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByIDAndUserID(db, req.ArticleID, userID, []string{"id", "status", "title", "slug", "category_id"}, []string{"User", "Category", "Tags"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID),
				zap.Uint("userID", userID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	updateFields := map[string]interface{}{
		"status": req.Body.Status,
	}

	if err := s.articleDAO.Update(db, article, updateFields); err != nil {
		logger.Error("[ArticleService] failed to update article status",
			zap.Uint("articleID", req.ArticleID),
			zap.String("status", string(req.Body.Status)),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(ctx context.Context, req *dto.DeleteArticleRequest) (rsp *dto.EmptyResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil {
		logger.Error("[ArticleService] request is nil")
		return nil, protocol.ErrBadRequest
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.EmptyResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != userID {
		logger.Error("[ArticleService] no permission to delete article",
			zap.Uint("articleID", article.ID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.articleDAO.Delete(db, article); err != nil {
		logger.Error("[ArticleService] failed to delete article",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListArticles 列出文章
func (s *articleService) ListArticles(ctx context.Context, req *dto.ListArticleRequest) (rsp *dto.ListArticleResponse, err error) {
	rsp = &dto.ListArticleResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.PageParam.Page,
			PageSize: req.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.QueryParam.Query,
			QueryFields: []string{"title", "slug"},
		},
	}

	articles, pageInfo, err := s.articleDAO.Paginate(db,
		[]string{
			"id", "slug", "title", "status", "user_id", "category_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"User", "Category", "Tags", "Comments"},
		param,
	)
	if err != nil {
		logger.Error("[ArticleService] failed to list articles", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *dto.Article {
		return s.buildArticleDTO(&article)
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

func (s *articleService) buildArticleDTO(article *model.Article) *dto.Article {
	return &dto.Article{
		ArticleID: article.ID,
		Title:     article.Title,
		Slug:      article.Slug,
		Status:    string(article.Status),
		User: &dto.User{
			UserID: article.User.ID,
			Name:   article.User.Name,
			Avatar: article.User.Avatar,
		},
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *dto.Tag {
			return &dto.Tag{
				TagID: tag.ID,
				Name:  tag.Name,
				Slug:  tag.Slug,
			}
		}),
		Comments: len(article.Comments),
	}
}
