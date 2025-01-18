package service

import (
	"errors"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
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
	CreateArticle(req *protocol.CreateArticleRequest) (rsp *protocol.CreateArticleResponse, err error)
	GetArticleInfo(req *protocol.GetArticleInfoRequest) (rsp *protocol.GetArticleInfoResponse, err error)
	UpdateArticle(req *protocol.UpdateArticleRequest) (rsp *protocol.UpdateArticleResponse, err error)
	UpdateArticleStatus(req *protocol.UpdateArticleStatusRequest) (rsp *protocol.UpdateArticleStatusResponse, err error)
	DeleteArticle(req *protocol.DeleteArticleRequest) (rsp *protocol.DeleteArticleResponse, err error)
	ListArticles(req *protocol.ListArticlesRequest) (rsp *protocol.ListArticlesResponse, err error)
}

type articleService struct {
	db                *gorm.DB
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
		db:                database.GetDBInstance(),
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
func (s *articleService) CreateArticle(req *protocol.CreateArticleRequest) (rsp *protocol.CreateArticleResponse, err error) {
	rsp = &protocol.CreateArticleResponse{}

	tags := []model.Tag{}
	tagChan, errChan := make(chan *model.Tag, len(req.Tags)), make(chan error, len(req.Tags))

	var wg sync.WaitGroup
	wg.Add(len(req.Tags))

	getTagFunc := func(tagSlug string) {
		defer wg.Done()
		tag, err := s.tagDAO.GetBySlug(s.db, tagSlug, []string{"id", "slug"}, []string{})
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
			logger.Logger.Error("[ArticleService] tag not found", zap.Error(err))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get tag", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	for tag := range tagChan {
		tags = append(tags, *tag)
	}

	category, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleService] category not found",
				zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
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

	if err := s.articleDAO.Create(s.db, article); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Logger.Error("[ArticleService] article slug duplicated",
				zap.Uint("userID", article.UserID),
				zap.String("title", article.Title),
				zap.String("slug", article.Slug))
			return nil, protocol.ErrDataExists
		}
		logger.Logger.Error("[ArticleService] failed to create article",
			zap.Uint("userID", article.UserID),
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
		UserID:      article.UserID,
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags:        lo.Map(article.Tags, func(tag model.Tag, _ int) string { return tag.Slug }),
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
func (s *articleService) GetArticleInfo(req *protocol.GetArticleInfoRequest) (rsp *protocol.GetArticleInfoResponse, err error) {
	rsp = &protocol.GetArticleInfoResponse{}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{
		"id", "slug", "title", "status", "user_id",
		"created_at", "updated_at", "published_at",
		"likes", "views",
	}, []string{"Comments", "Tags"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	// 如果文章不是公开的，则只有作者本人可以查看
	if article.UserID != req.UserID && article.Status != model.ArticleStatusPublish {
		logger.Logger.Error("[ArticleService] no permission to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Article = &protocol.Article{
		ArticleID:   article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      string(article.Status),
		UserID:      article.UserID,
		CreatedAt:   article.CreatedAt.Format(time.DateTime),
		UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
		PublishedAt: article.PublishedAt.Format(time.DateTime),
		Likes:       article.Likes,
		Views:       article.Views,
		Tags:        lo.Map(article.Tags, func(tag model.Tag, _ int) string { return tag.Slug }),
		Comments:    len(article.Comments),
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
func (s *articleService) UpdateArticle(req *protocol.UpdateArticleRequest) (rsp *protocol.UpdateArticleResponse, err error) {
	rsp = &protocol.UpdateArticleResponse{}

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
		logger.Logger.Warn("[ArticleService] no fields to update",
			zap.Uint("userID", req.UserID),
			zap.Uint("articleID", req.ArticleID),
			zap.Any("updateFields", updateFields))
		return rsp, nil
	}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{"id", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != req.UserID {
		logger.Logger.Error("[ArticleService] no permission to update article",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.articleDAO.Update(s.db, article, updateFields); err != nil {
		logger.Logger.Error("[ArticleService] failed to update article",
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
func (s *articleService) UpdateArticleStatus(req *protocol.UpdateArticleStatusRequest) (rsp *protocol.UpdateArticleStatusResponse, err error) {
	rsp = &protocol.UpdateArticleStatusResponse{}

	article, err := s.articleDAO.GetByIDAndUserID(s.db, req.ArticleID, req.UserID, []string{"id", "status", "title", "slug", "category_id"}, []string{"User", "Category", "Tags"})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.Status == req.Status {
		logger.Logger.Warn("[ArticleService] article status not changed",
			zap.Uint("articleID", req.ArticleID),
			zap.String("status", string(req.Status)))
		return rsp, nil
	}

	if req.Status == model.ArticleStatusPublish {
		if err := s.articleDAO.Update(s.db, article, map[string]interface{}{
			"status":       req.Status,
			"published_at": time.Now(),
		}); err != nil {
			logger.Logger.Error("[ArticleService] failed to update article status",
				zap.Uint("articleID", article.ID),
				zap.String("status", string(req.Status)),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else if req.Status == model.ArticleStatusDraft {
		if err := s.articleDAO.Update(s.db, article, map[string]interface{}{
			"status":       req.Status,
			"published_at": nil,
		}); err != nil {
			logger.Logger.Error("[ArticleService] failed to update article status",
				zap.Uint("articleID", article.ID),
				zap.String("status", string(req.Status)),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
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
func (s *articleService) DeleteArticle(req *protocol.DeleteArticleRequest) (rsp *protocol.DeleteArticleResponse, err error) {
	rsp = &protocol.DeleteArticleResponse{}

	article, err := s.articleDAO.GetByIDAndUserID(s.db, req.ArticleID, req.UserID, []string{"id", "slug"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if err := s.articleDAO.Delete(s.db, article); err != nil {
		logger.Logger.Error("[ArticleService] failed to delete article",
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
func (s *articleService) ListArticles(req *protocol.ListArticlesRequest) (rsp *protocol.ListArticlesResponse, err error) {
	rsp = &protocol.ListArticlesResponse{}

	articles, pageInfo, err := s.articleDAO.PaginateByStatus(
		s.db,
		model.ArticleStatusPublish,
		[]string{
			"id", "slug", "title", "status", "user_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"Comments", "Tags"},
		req.PageParam.Page, req.PageParam.PageSize,
	)
	if err != nil {
		logger.Logger.Error("[ArticleService] failed to paginate articles", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *protocol.Article {
		return &protocol.Article{
			ArticleID:   article.ID,
			Title:       article.Title,
			Slug:        article.Slug,
			Status:      string(article.Status),
			UserID:      article.UserID,
			CreatedAt:   article.CreatedAt.Format(time.DateTime),
			UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
			PublishedAt: article.PublishedAt.Format(time.DateTime),
			Likes:       article.Likes,
			Views:       article.Views,
			Tags:        lo.Map(article.Tags, func(tag model.Tag, _ int) string { return tag.Slug }),
			Comments:    len(article.Comments),
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
