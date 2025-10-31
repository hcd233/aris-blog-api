package service

import (
	"context"
	"errors"
	"time"

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

// CategoryService 分类服务
type CategoryService interface {
	CreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (rsp *dto.CategoryCreateResponse, err error)
	GetCategoryInfo(ctx context.Context, req *dto.CategoryGetRequest) (rsp *dto.CategoryGetResponse, err error)
	GetRootCategory(ctx context.Context, req *dto.CategoryGetRootRequest) (rsp *dto.CategoryGetRootResponse, err error)
	UpdateCategory(ctx context.Context, req *dto.CategoryUpdateRequest) (rsp *dto.CategoryUpdateResponse, err error)
	DeleteCategory(ctx context.Context, req *dto.CategoryDeleteRequest) (rsp *dto.CategoryDeleteResponse, err error)
	ListChildrenCategories(ctx context.Context, req *dto.CategoryListChildrenCategoriesRequest) (rsp *dto.CategoryListChildrenCategoriesResponse, err error)
	ListChildrenArticles(ctx context.Context, req *dto.CategoryListChildrenArticlesRequest) (rsp *dto.CategoryListChildrenArticlesResponse, err error)
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
func (s *categoryService) CreateCategory(ctx context.Context, req *dto.CategoryCreateRequest) (rsp *dto.CategoryCreateResponse, err error) {
	rsp = &dto.CategoryCreateResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req.Body == nil {
		return nil, protocol.ErrBadRequest
	}

	name := req.Body.Name
	parentID := req.Body.ParentID

	var parentCategory *model.Category
	if parentID == 0 {
		parentCategory, err = s.categoryDAO.GetRootByUserID(db, req.UserID, []string{"id"}, []string{})
	} else {
		parentCategory, err = s.categoryDAO.GetByID(db, parentID, []string{"id"}, []string{})
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] parent category not found",
				zap.Uint("parentID", parentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get parent category",
			zap.Uint("parentID", parentID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	category := &model.Category{
		Name:     name,
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

	rsp.Category = &dto.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetCategoryInfo 获取分类信息
func (s *categoryService) GetCategoryInfo(ctx context.Context, req *dto.CategoryGetRequest) (rsp *dto.CategoryGetResponse, err error) {
	rsp = &dto.CategoryGetResponse{}

	logger := logger.WithCtx(ctx)
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
			zap.Uint("categoryUserID", category.UserID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Category = &dto.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetRootCategory 获取根分类
func (s *categoryService) GetRootCategory(ctx context.Context, req *dto.CategoryGetRootRequest) (rsp *dto.CategoryGetRootResponse, err error) {
	rsp = &dto.CategoryGetRootResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	rootCategory, err := s.categoryDAO.GetRootByUserID(db, req.UserID, []string{"id", "name", "user_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CategoryService] root category not found")
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CategoryService] failed to get root category", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Category = &dto.Category{
		CategoryID: rootCategory.ID,
		Name:       rootCategory.Name,
		ParentID:   rootCategory.ParentID,
		CreatedAt:  rootCategory.CreatedAt.Format(time.DateTime),
		UpdatedAt:  rootCategory.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(ctx context.Context, req *dto.CategoryUpdateRequest) (rsp *dto.CategoryUpdateResponse, err error) {
	rsp = &dto.CategoryUpdateResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req.Body == nil {
		return nil, protocol.ErrBadRequest
	}

	updateFields := make(map[string]interface{})
	if req.Body.Name != "" {
		updateFields["name"] = req.Body.Name
	}
	if req.Body.ParentID != 0 {
		updateFields["parent_id"] = req.Body.ParentID
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

	rsp.Category = &dto.Category{
		CategoryID: category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		CreatedAt:  category.CreatedAt.Format(time.DateTime),
		UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(ctx context.Context, req *dto.CategoryDeleteRequest) (rsp *dto.CategoryDeleteResponse, err error) {
	rsp = &dto.CategoryDeleteResponse{}

	logger := logger.WithCtx(ctx)
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
func (s *categoryService) ListChildrenCategories(ctx context.Context, req *dto.CategoryListChildrenCategoriesRequest) (rsp *dto.CategoryListChildrenCategoriesResponse, err error) {
	rsp = &dto.CategoryListChildrenCategoriesResponse{}

	logger := logger.WithCtx(ctx)
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

	paginate := req.PaginationQuery.ToPaginateParam()
	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     paginate.PageParam.Page,
			PageSize: paginate.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       paginate.QueryParam.Query,
			QueryFields: []string{"name"},
		},
	}
	categories, pageInfo, err := s.categoryDAO.PaginateChildren(db, parentCategory,
		[]string{"id", "name", "parent_id", "created_at", "updated_at"}, []string{},
		param)
	if err != nil {
		logger.Error("[CategoryService] failed to paginate children categories",
			zap.Uint("parentID", parentCategory.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.UserID != parentCategory.UserID {
		logger.Error("[CategoryService] no permission to list children categories",
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	rsp.Categories = lo.Map(*categories, func(category model.Category, _ int) *dto.Category {
		return &dto.Category{
			CategoryID: category.ID,
			Name:       category.Name,
			ParentID:   category.ParentID,
			CreatedAt:  category.CreatedAt.Format(time.DateTime),
			UpdatedAt:  category.UpdatedAt.Format(time.DateTime),
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// ListChildrenArticles 列出子文章
func (s *categoryService) ListChildrenArticles(ctx context.Context, req *dto.CategoryListChildrenArticlesRequest) (rsp *dto.CategoryListChildrenArticlesResponse, err error) {
	rsp = &dto.CategoryListChildrenArticlesResponse{}

	logger := logger.WithCtx(ctx)
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
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	paginate := req.PaginationQuery.ToPaginateParam()
	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     paginate.PageParam.Page,
			PageSize: paginate.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       paginate.QueryParam.Query,
			QueryFields: []string{"title", "slug"},
		},
	}
	articles, pageInfo, err := s.articleDAO.PaginateByCategoryID(db, parentCategory.ID,
		[]string{
			"id", "slug", "title", "status", "user_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"Tags", "Comments"},
		param)
	if err != nil {
		logger.Error("[CategoryService] failed to paginate children articles",
			zap.Uint("categoryID", parentCategory.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *dto.Article {
		return &dto.Article{
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
			Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *dto.Tag {
				return &dto.Tag{
					TagID: tag.ID,
					Name:  tag.Name,
					Slug:  tag.Slug,
				}
			}),
			Comments: int(len(article.Comments)),
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
