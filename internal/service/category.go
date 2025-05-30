package service

import (
	"context"
	"errors"
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

// CategoryService 分类服务
type CategoryService interface {
	CreateCategory(ctx context.Context, req *protocol.CreateCategoryRequest) (rsp *protocol.CreateCategoryResponse, err error)
	GetCategoryInfo(ctx context.Context, req *protocol.GetCategoryInfoRequest) (rsp *protocol.GetCategoryInfoResponse, err error)
	GetRootCategory(ctx context.Context, req *protocol.GetRootCategoryRequest) (rsp *protocol.GetRootCategoryResponse, err error)
	UpdateCategory(ctx context.Context, req *protocol.UpdateCategoryRequest) (rsp *protocol.UpdateCategoryResponse, err error)
	DeleteCategory(ctx context.Context, req *protocol.DeleteCategoryRequest) (rsp *protocol.DeleteCategoryResponse, err error)
	ListChildrenCategories(ctx context.Context, req *protocol.ListChildrenCategoriesRequest) (rsp *protocol.ListChildrenCategoriesResponse, err error)
	ListChildrenArticles(ctx context.Context, req *protocol.ListChildrenArticlesRequest) (rsp *protocol.ListChildrenArticlesResponse, err error)
}

type categoryService struct {
	userDAO     *dao.UserDAO
	categoryDAO *dao.CategoryDAO
	articleDAO  *dao.ArticleDAO
}

// NewCategoryService 创建分类服务
func NewCategoryService() CategoryService {
	return &categoryService{
		userDAO:     dao.GetUserDAO(),
		categoryDAO: dao.GetCategoryDAO(),
		articleDAO:  dao.GetArticleDAO(),
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(ctx context.Context, req *protocol.CreateCategoryRequest) (rsp *protocol.CreateCategoryResponse, err error) {
	rsp = &protocol.CreateCategoryResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	var parentCategory *model.Category
	if req.ParentID == 0 {
		parentCategory, err = s.categoryDAO.GetRootByUserID(db, req.UserID, []string{"id"}, []string{})
	} else {
		parentCategory, err = s.categoryDAO.GetByID(db, req.ParentID, []string{"id"}, []string{})
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] parent category not found",
				zap.Uint("parentID", req.ParentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get parent category",
			zap.Uint("parentID", req.ParentID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	category := &model.Category{
		Name:     req.Name,
		ParentID: parentCategory.ID,
		UserID:   req.UserID,
	}

	if err := s.categoryDAO.Create(db, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Error("[CategoryService] duplicated category",
				zap.String("name", category.Name),
				zap.Uint("parentID", category.ParentID))
			return nil, protocol.ErrDataExists
		}
		logger.Error("[CategoryService] failed to create category",
			zap.String("name", category.Name),
			zap.Uint("parentID", category.ParentID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Category = &protocol.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetCategoryInfo 获取分类信息
func (s *categoryService) GetCategoryInfo(ctx context.Context, req *protocol.GetCategoryInfoRequest) (rsp *protocol.GetCategoryInfoResponse, err error) {
	rsp = &protocol.GetCategoryInfoResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	category, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id", "name", "user_id", "parent_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != category.UserID {
		logger.Error("[CategoryService] no permission to get category",
			zap.Uint("userID", req.UserID),
			zap.Uint("categoryUserID", category.UserID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Category = &protocol.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetRootCategory 获取根分类
func (s *categoryService) GetRootCategory(ctx context.Context, req *protocol.GetRootCategoryRequest) (rsp *protocol.GetRootCategoryResponse, err error) {
	rsp = &protocol.GetRootCategoryResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	rootCategory, err := s.categoryDAO.GetRootByUserID(db, req.UserID, []string{"id", "name", "user_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] root category not found", zap.Uint("userID", req.UserID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get root category", zap.Uint("userID", req.UserID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Category = &protocol.Category{
		CategoryID: rootCategory.ID,
		Name:       rootCategory.Name,
		ParentID:   rootCategory.ParentID,
		CreatedAt:  rootCategory.CreatedAt.Format(time.DateTime),
		UpdatedAt:  rootCategory.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(ctx context.Context, req *protocol.UpdateCategoryRequest) (rsp *protocol.UpdateCategoryResponse, err error) {
	rsp = &protocol.UpdateCategoryResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	updateFields := make(map[string]interface{})
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.ParentID != 0 {
		updateFields["parent_id"] = req.ParentID
	}

	if len(updateFields) == 0 {
		logger.Warn("[CategoryService] no fields to update",
			zap.Uint("categoryID", req.CategoryID))
		return rsp, nil
	}

	category, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id", "name", "user_id", "parent_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != category.UserID {
		logger.Error("[CategoryService] no permission to update category",
			zap.Uint("userID", req.UserID),
			zap.Uint("categoryUserID", category.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.categoryDAO.Update(db, category, updateFields); err != nil {
		logger.Error("[CategoryService] failed to update category",
			zap.Uint("categoryID", category.ID),
			zap.Any("updateFields", updateFields),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Category = &protocol.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(ctx context.Context, req *protocol.DeleteCategoryRequest) (rsp *protocol.DeleteCategoryResponse, err error) {
	rsp = &protocol.DeleteCategoryResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	category, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id", "name", "user_id", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != category.UserID {
		logger.Error("[CategoryService] no permission to delete category",
			zap.Uint("userID", req.UserID),
			zap.Uint("categoryUserID", category.UserID))
		return nil, protocol.ErrNoPermission
	}

	if category.ParentID == 0 {
		logger.Error("[CategoryService] root category cannot be deleted",
			zap.Uint("categoryID", category.ID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.categoryDAO.DeleteReclusiveByID(db, category.ID, []string{"id", "name"}, []string{}); err != nil {
		logger.Error("[CategoryService] failed to delete category",
			zap.Uint("categoryID", category.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListChildrenCategories 列出子分类
func (s *categoryService) ListChildrenCategories(ctx context.Context, req *protocol.ListChildrenCategoriesRequest) (rsp *protocol.ListChildrenCategoriesResponse, err error) {
	rsp = &protocol.ListChildrenCategoriesResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	parentCategory, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id", "name", "user_id", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] parent category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get parent category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	categories, pageInfo, err := s.categoryDAO.PaginateChildren(db, parentCategory,
		[]string{"id", "name", "parent_id", "created_at", "updated_at"}, []string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[CategoryService] failed to paginate children categories",
			zap.Uint("parentID", parentCategory.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != parentCategory.UserID {
		logger.Error("[CategoryService] no permission to list children categories",
			zap.Uint("userID", req.UserID),
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Categories = lo.Map(*categories, func(category model.Category, _ int) *protocol.Category {
		return &protocol.Category{
			CategoryID: category.ID,
			Name:       category.Name,
			ParentID:   category.ParentID,
			CreatedAt:  category.CreatedAt.Format(time.DateTime),
			UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// ListChildrenArticles 列出子文章
func (s *categoryService) ListChildrenArticles(ctx context.Context, req *protocol.ListChildrenArticlesRequest) (rsp *protocol.ListChildrenArticlesResponse, err error) {
	rsp = &protocol.ListChildrenArticlesResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	parentCategory, err := s.categoryDAO.GetByID(db, req.CategoryID, []string{"id", "name", "user_id", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] parent category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get parent category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != parentCategory.UserID {
		logger.Error("[CategoryService] no permission to list children articles",
			zap.Uint("userID", req.UserID),
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	articles, pageInfo, err := s.articleDAO.PaginateByCategoryID(db, parentCategory.ID,
		[]string{
			"id", "slug", "title", "status", "user_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"Tags", "Comments"},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[CategoryService] failed to paginate children articles",
			zap.Uint("categoryID", parentCategory.ID),
			zap.Error(err))
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
			Comments:    int(len(article.Comments)),
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
