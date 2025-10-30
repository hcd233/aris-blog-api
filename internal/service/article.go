package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleService 文章服务
//
//	author centonhuang
//	update 2025-01-05 15:23:26
type ArticleService interface {
	CreateArticle(ctx context.Context, req *protocol.CreateArticleRequest) (rsp *protocol.CreateArticleResponse, err error)
	GetArticleInfo(ctx context.Context, req *protocol.GetArticleInfoRequest) (rsp *protocol.GetArticleInfoResponse, err error)
	GetArticleInfoBySlug(ctx context.Context, req *protocol.GetArticleInfoBySlugRequest) (rsp *protocol.GetArticleInfoBySlugResponse, err error)
	UpdateArticle(ctx context.Context, req *protocol.UpdateArticleRequest) (rsp *protocol.UpdateArticleResponse, err error)
	UpdateArticleStatus(ctx context.Context, req *protocol.UpdateArticleStatusRequest) (rsp *protocol.UpdateArticleStatusResponse, err error)
	DeleteArticle(ctx context.Context, req *protocol.DeleteArticleRequest) (rsp *protocol.DeleteArticleResponse, err error)
	ListArticles(ctx context.Context, req *protocol.ListArticlesRequest) (rsp *protocol.ListArticlesResponse, err error)
}

type articleService struct {
	userDAO           *dao.UserDAO
	tagDAO            *dao.TagDAO
	categoryDAO       *dao.CategoryDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
}

// NewArticleService 创建文章服务
//
//	return ArticleService
//	author centonhuang
//	update 2025-01-05 15:23:26
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
//
//	receiver s *articleService
//	param req *protocol.CreateArticleRequest
//	return rsp *protocol.CreateArticleResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) CreateArticle(ctx context.Context, req *protocol.CreateArticleRequest) (rsp *protocol.CreateArticleResponse, err error) {
	rsp = &protocol.CreateArticleResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	tags := []model.Tag{}
	tagChan, errChan := make(chan *model.Tag, len(req.Tags)), make(chan error, len(req.Tags))

	var wg sync.WaitGroup
	wg.Add(len(req.Tags))

	getTagFunc := func(tagSlug string) {
		defer wg.Done()
		tag, err := s.tagDAO.GetBySlug(db, tagSlug, []string{"id", "slug"}, []string{})
		if err != nil {
			errChan <- err
			return
		}
		tagChan <- tag
	}

	for _, tagSlug := range req.Tags {
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

	category, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleService] category not found",
				zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article := &model.Article{
		UserID:   req.UserID,
		Status:   model.ArticleStatusDraft,
		Title:    req.Title,
		Slug:     req.Slug,
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

	rsp.Article = &protocol.Article{
		ArticleID:   article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      string(article.Status),
		User:        nil,
		Category:    nil,
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
//
//	receiver s *articleService
//	param req *protocol.GetArticleInfoRequest
//	return rsp *protocol.GetArticleInfoResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) GetArticleInfo(ctx context.Context, req *protocol.GetArticleInfoRequest) (rsp *protocol.GetArticleInfoResponse, err error) {
	rsp = &protocol.GetArticleInfoResponse{}

	logger := logger.WithCtx(ctx)
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

	// 如果文章不是公开的，则只有作者本人可以查看
	if article.UserID != req.UserID && article.Status != model.ArticleStatusPublish {
		logger.Error("[ArticleService] no permission to get article",
			zap.Uint("articleID", req.ArticleID),
		)
		return nil, protocol.ErrNoPermission
	}

	rsp.Article = &protocol.Article{
		ArticleID: article.ID,
		Title:     article.Title,
		Slug:      article.Slug,
		Status:    string(article.Status),
		User: &dto.User{
			UserID: article.User.ID,
			Name:   article.User.Name,
			Avatar: article.User.Avatar,
		},
		Category: &protocol.Category{
			CategoryID: article.CategoryID,
			Name:       article.Category.Name,
		},
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *protocol.Tag {
			return &protocol.Tag{
				TagID: tag.ID,
				Name:  tag.Name,
				Slug:  tag.Slug,
			}
		}),
		Comments: len(article.Comments),
	}

	return rsp, nil
}

// GetArticleInfoBySlug 获取文章信息
//
//	receiver s *articleService
//	param req *protocol.GetArticleInfoBySlugRequest
//	return rsp *protocol.GetArticleInfoBySlugResponse
//	return err error
//	author centonhuang
//	update 2025-01-19 15:23:26
func (s *articleService) GetArticleInfoBySlug(ctx context.Context, req *protocol.GetArticleInfoBySlugRequest) (rsp *protocol.GetArticleInfoBySlugResponse, err error) {
	rsp = &protocol.GetArticleInfoBySlugResponse{}

	logger := logger.WithCtx(ctx)
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

	// 如果文章不是公开的，则只有作者本人可以查看
	if article.UserID != user.ID && article.Status != model.ArticleStatusPublish {
		logger.Error("[ArticleService] no permission to get article",
			zap.String("ArticleSlug", req.ArticleSlug),
			zap.String("AuthorName", req.AuthorName),
		)
		return nil, protocol.ErrNoPermission
	}

	rsp.Article = &protocol.Article{
		ArticleID: article.ID,
		Title:     article.Title,
		Slug:      article.Slug,
		Status:    string(article.Status),
		User: &dto.User{
			UserID: article.User.ID,
			Name:   article.User.Name,
			Avatar: article.User.Avatar,
		},
		Category: &protocol.Category{
			CategoryID: article.CategoryID,
			Name:       article.Category.Name,
		},
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *protocol.Tag {
			return &protocol.Tag{
				TagID: tag.ID,
				Name:  tag.Name,
				Slug:  tag.Slug,
			}
		}),
		Comments: len(article.Comments),
	}

	return rsp, nil
}

// UpdateArticle 更新文章
//
//	receiver s *articleService
//	param req *protocol.UpdateArticleRequest
//	return rsp *protocol.UpdateArticleResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) UpdateArticle(ctx context.Context, req *protocol.UpdateArticleRequest) (rsp *protocol.UpdateArticleResponse, err error) {
	rsp = &protocol.UpdateArticleResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	updateFields := make(map[string]interface{})
	if req.UpdatedTitle != "" {
		updateFields["title"] = req.UpdatedTitle
	}
	if req.UpdatedSlug != "" {
		updateFields["slug"] = req.UpdatedSlug
	}
	if req.UpdatedCategoryID != 0 {
		updateFields["category_id"] = req.UpdatedCategoryID
	}

	if len(updateFields) == 0 {
		logger.Warn("[ArticleService] no fields to update",
			zap.Uint("articleID", req.ArticleID),
			zap.Any("updateFields", updateFields))
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

	if article.UserID != req.UserID {
		logger.Error("[ArticleService] no permission to update article",
			zap.Uint("articleID", req.ArticleID),
		)
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
//
//	receiver s *articleService
//	param req *protocol.UpdateArticleStatusRequest
//	return rsp *protocol.UpdateArticleStatusResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) UpdateArticleStatus(ctx context.Context, req *protocol.UpdateArticleStatusRequest) (rsp *protocol.UpdateArticleStatusResponse, err error) {
	rsp = &protocol.UpdateArticleStatusResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByIDAndUserID(db, req.ArticleID, req.UserID, []string{"id", "status", "title", "slug", "category_id"}, []string{"User", "Category", "Tags"})
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

	if article.Status == req.Status {
		logger.Warn("[ArticleService] article status not changed",
			zap.Uint("articleID", req.ArticleID),
			zap.String("status", string(req.Status)))
		return rsp, nil
	}

	updateFields := map[string]interface{}{
		"status": req.Status,
	}
	switch req.Status {
	case model.ArticleStatusPublish:
		updateFields["published_at"] = time.Now().UTC()

	case model.ArticleStatusDraft:
		updateFields["published_at"] = nil
	}

	if err := s.articleDAO.Update(db, article, updateFields); err != nil {
		logger.Error("[ArticleService] failed to update article status",
			zap.Uint("articleID", article.ID),
			zap.String("status", string(req.Status)),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}
	return rsp, nil
}

// DeleteArticle 删除文章
//
//	receiver s *articleService
//	param req *protocol.DeleteArticleRequest
//	return rsp *protocol.DeleteArticleResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) DeleteArticle(ctx context.Context, req *protocol.DeleteArticleRequest) (rsp *protocol.DeleteArticleResponse, err error) {
	rsp = &protocol.DeleteArticleResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByIDAndUserID(db, req.ArticleID, req.UserID, []string{"id", "slug"}, []string{})
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

	if err := s.articleDAO.Delete(db, article); err != nil {
		logger.Error("[ArticleService] failed to delete article",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListArticles 列出文章
//
//	receiver s *articleService
//	param req *protocol.ListArticlesRequest
//	return rsp *protocol.ListArticlesResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 15:23:26
func (s *articleService) ListArticles(ctx context.Context, req *protocol.ListArticlesRequest) (rsp *protocol.ListArticlesResponse, err error) {
	rsp = &protocol.ListArticlesResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     req.PaginateParam.Page,
			PageSize: req.PaginateParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.PaginateParam.Query,
			QueryFields: []string{"title"},
		},
	}
	articles, pageInfo, err := s.articleDAO.Paginate(
		db,
		[]string{
			"id", "slug", "title", "status", "user_id", "category_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"User", "Category", "Tags", "Comments"},
		param,
	)
	if err != nil {
		logger.Error("[ArticleService] failed to paginate articles", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *protocol.Article {
		return &protocol.Article{
			ArticleID: article.ID,
			Title:     article.Title,
			Slug:      article.Slug,
			Status:    string(article.Status),
			User: &dto.User{
				UserID: article.User.ID,
				Name:   article.User.Name,
				Avatar: article.User.Avatar,
			},
			Category: &protocol.Category{
				CategoryID: article.CategoryID,
				Name:       article.Category.Name,
			},
			CreatedAt:   article.CreatedAt.Format(time.DateTime),
			UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
			PublishedAt: article.PublishedAt.Format(time.DateTime),
			Likes:       article.Likes,
			Views:       article.Views,
			Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *protocol.Tag {
				return &protocol.Tag{
					TagID: tag.ID,
					Name:  tag.Name,
					Slug:  tag.Slug,
				}
			}),
			Comments: len(article.Comments),
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
