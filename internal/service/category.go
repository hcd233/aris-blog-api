package service

import (
	"errors"
	"time"

	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CategoryService 分类服务
type CategoryService interface {
	CreateCategory(req *protocol.CreateCategoryRequest) (rsp *protocol.CreateCategoryResponse, err error)
	GetCategoryInfo(req *protocol.GetCategoryInfoRequest) (rsp *protocol.GetCategoryInfoResponse, err error)
	GetRootCategory(req *protocol.GetRootCategoryRequest) (rsp *protocol.GetRootCategoryResponse, err error)
	UpdateCategory(req *protocol.UpdateCategoryRequest) (rsp *protocol.UpdateCategoryResponse, err error)
	DeleteCategory(req *protocol.DeleteCategoryRequest) (rsp *protocol.DeleteCategoryResponse, err error)
	ListChildrenCategories(req *protocol.ListChildrenCategoriesRequest) (rsp *protocol.ListChildrenCategoriesResponse, err error)
	ListChildrenArticles(req *protocol.ListChildrenArticlesRequest) (rsp *protocol.ListChildrenArticlesResponse, err error)
}

type categoryService struct {
	db          *gorm.DB
	userDAO     *dao.UserDAO
	categoryDAO *dao.CategoryDAO
	articleDAO  *dao.ArticleDAO
}

// NewCategoryService 创建分类服务
func NewCategoryService() CategoryService {
	return &categoryService{
		db:          database.GetDBInstance(),
		userDAO:     dao.GetUserDAO(),
		categoryDAO: dao.GetCategoryDAO(),
		articleDAO:  dao.GetArticleDAO(),
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(req *protocol.CreateCategoryRequest) (rsp *protocol.CreateCategoryResponse, err error) {
	rsp = &protocol.CreateCategoryResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to create category",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	var parentCategory *model.Category
	if req.ParentID == 0 {
		parentCategory, err = s.categoryDAO.GetRootByUserID(s.db, user.ID, []string{"id"}, []string{})
	} else {
		parentCategory, err = s.categoryDAO.GetByID(s.db, req.ParentID, []string{"id"}, []string{})
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] parent category not found",
				zap.Uint("parentID", req.ParentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get parent category",
			zap.Uint("parentID", req.ParentID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	category := &model.Category{
		Name:     req.Name,
		ParentID: parentCategory.ID,
		UserID:   user.ID,
	}

	if err := s.categoryDAO.Create(s.db, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Logger.Error("[CategoryService] duplicated category",
				zap.String("name", category.Name),
				zap.Uint("parentID", category.ParentID))
			return nil, protocol.ErrDataExists
		}
		logger.Logger.Error("[CategoryService] failed to create category",
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
func (s *categoryService) GetCategoryInfo(req *protocol.GetCategoryInfoRequest) (rsp *protocol.GetCategoryInfoResponse, err error) {
	rsp = &protocol.GetCategoryInfoResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to get category",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	_, err = s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	category, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id", "name", "parent_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
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

// GetRootCategory 获取根分类
func (s *categoryService) GetRootCategory(req *protocol.GetRootCategoryRequest) (rsp *protocol.GetRootCategoryResponse, err error) {
	rsp = &protocol.GetRootCategoryResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to get root category",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rootCategory, err := s.categoryDAO.GetRootByUserID(s.db, user.ID, []string{"id", "name", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] root category not found", zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get root category", zap.Uint("userID", user.ID), zap.Error(err))
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
func (s *categoryService) UpdateCategory(req *protocol.UpdateCategoryRequest) (rsp *protocol.UpdateCategoryResponse, err error) {
	rsp = &protocol.UpdateCategoryResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to update category",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	_, err = s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	updateFields := make(map[string]interface{})
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.ParentID != 0 {
		updateFields["parent_id"] = req.ParentID
	}

	if len(updateFields) == 0 {
		logger.Logger.Warn("[CategoryService] no fields to update",
			zap.String("userName", req.UserName),
			zap.Uint("categoryID", req.CategoryID))
		return rsp, nil
	}

	category, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id", "name", "parent_id", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if err := s.categoryDAO.Update(s.db, category, updateFields); err != nil {
		logger.Logger.Error("[CategoryService] failed to update category",
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
func (s *categoryService) DeleteCategory(req *protocol.DeleteCategoryRequest) (rsp *protocol.DeleteCategoryResponse, err error) {
	rsp = &protocol.DeleteCategoryResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to delete category",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	category, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != category.UserID {
		logger.Logger.Error("[CategoryService] no permission to delete category",
			zap.Uint("userID", user.ID),
			zap.Uint("categoryUserID", category.UserID))
		return nil, protocol.ErrNoPermission
	}

	if category.ParentID == 0 {
		logger.Logger.Error("[CategoryService] root category cannot be deleted",
			zap.Uint("categoryID", category.ID))
		return nil, protocol.ErrInternalError
	}

	if err := s.categoryDAO.DeleteReclusiveByID(s.db, category.ID, []string{"id", "name"}, []string{}); err != nil {
		logger.Logger.Error("[CategoryService] failed to delete category",
			zap.Uint("categoryID", category.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListChildrenCategories 列出子分类
func (s *categoryService) ListChildrenCategories(req *protocol.ListChildrenCategoriesRequest) (rsp *protocol.ListChildrenCategoriesResponse, err error) {
	rsp = &protocol.ListChildrenCategoriesResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to list children categories",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	parentCategory, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] parent category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get parent category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != parentCategory.UserID {
		logger.Logger.Info("[CategoryService] no permission to list children categories",
			zap.Uint("userID", user.ID),
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	categories, pageInfo, err := s.categoryDAO.PaginateChildren(s.db, parentCategory,
		[]string{"id", "name", "parent_id", "created_at", "updated_at"}, []string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[CategoryService] failed to paginate children categories",
			zap.Uint("parentID", parentCategory.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
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
func (s *categoryService) ListChildrenArticles(req *protocol.ListChildrenArticlesRequest) (rsp *protocol.ListChildrenArticlesResponse, err error) {
	rsp = &protocol.ListChildrenArticlesResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Error("[CategoryService] no permission to list children articles",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	parentCategory, err := s.categoryDAO.GetByID(s.db, req.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CategoryService] parent category not found", zap.Uint("categoryID", req.CategoryID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CategoryService] failed to get parent category", zap.Uint("categoryID", req.CategoryID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if user.ID != parentCategory.UserID {
		logger.Logger.Error("[CategoryService] no permission to list children articles",
			zap.Uint("userID", user.ID),
			zap.Uint("categoryUserID", parentCategory.UserID))
		return nil, protocol.ErrNoPermission
	}

	articles, pageInfo, err := s.articleDAO.PaginateByCategoryID(s.db, parentCategory.ID,
		[]string{"id", "title", "slug", "created_at", "updated_at"}, []string{"User", "Tags", "Comments"},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[CategoryService] failed to paginate children articles",
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
			Author:      article.User.Name,
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
