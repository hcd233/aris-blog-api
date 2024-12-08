package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// CategoryService 分类服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type CategoryService interface {
	CreateCategoryHandler(c *gin.Context)
	GetCategoryInfoHandler(c *gin.Context)
	UpdateCategoryInfoHandler(c *gin.Context)
	DeleteCategoryHandler(c *gin.Context)
	ListRootCategoriesHandler(c *gin.Context)
	ListChildrenCategoriesHandler(c *gin.Context)
	ListChildrenArticlesHandler(c *gin.Context)
}

type categoryService struct {
	db          *gorm.DB
	userDAO     *dao.UserDAO
	categoryDAO *dao.CategoryDAO
	articleDAO  *dao.ArticleDAO
}

// NewCategoryService 创建分类服务
//
//	@return CategoryService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewCategoryService() CategoryService {
	return &categoryService{
		db:          database.GetDBInstance(),
		userDAO:     dao.GetUserDAO(),
		categoryDAO: dao.GetCategoryDAO(),
		articleDAO:  dao.GetArticleDAO(),
	}
}

// CreateCategoryHandler 创建分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-28 07:03:28
func (s *categoryService) CreateCategoryHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateCategoryBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's category",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	var parentCategory *model.Category
	if body.ParentID == 0 {
		parentCategory, err = s.categoryDAO.GetRootByUserID(s.db, user.ID, []string{"id"}, []string{})
	} else {
		parentCategory, err = s.categoryDAO.GetByID(s.db, body.ParentID, []string{"id"}, []string{})
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}
	body.ParentID = parentCategory.ID

	category := &model.Category{
		Name:     body.Name,
		ParentID: body.ParentID,
		UserID:   user.ID,
	}

	if err := s.categoryDAO.Create(s.db, category); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: category.GetBasicInfo(),
	})
}

// GetCategoryInfoHandler 获取分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 04:58:27
func (s *categoryService) GetCategoryInfoHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	_, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	category, err := s.categoryDAO.GetByID(s.db, uri.CategoryID, []string{"id", "name", "parent_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"category": category.GetBasicInfo(),
		},
	})
}

// ListRootCategoriesHandler 获取根分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 03:56:26
func (s *categoryService) ListRootCategoriesHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's root category",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	rootCategory, err := s.categoryDAO.GetRootByUserID(s.db, user.ID, []string{"id", "name"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"category": rootCategory.GetBasicInfo(),
		},
	})
}

// UpdateCategoryInfoHandler 更新分类信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 03:45:55
func (s *categoryService) UpdateCategoryInfoHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	body := c.MustGet("body").(*protocol.UpdateCategoryBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's category",
		})
		return
	}

	_, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	updateFields := make(map[string]interface{})

	if body.Name != "" {
		updateFields["name"] = body.Name
	}

	if body.ParentID != 0 {
		updateFields["parent_id"] = body.ParentID
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateCategoryError,
			Message: "No fields to update",
		})
		return
	}

	category, err := s.categoryDAO.GetByID(s.db, uri.CategoryID, []string{"id", "name", "parent_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	err = s.categoryDAO.Update(s.db, category, updateFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"category": category.GetBasicInfo(),
		},
	})
}

// DeleteCategoryHandler 删除分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 04:55:08
func (s *categoryService) DeleteCategoryHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's category",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	category, err := s.categoryDAO.GetByID(s.db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	if user.ID != category.UserID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's category",
		})
		return
	}

	if category.ParentID == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteCategoryError,
			Message: "Root category can not be deleted",
		})
		return
	}

	err = s.categoryDAO.DeleteReclusiveByID(s.db, category.ID, []string{"id", "name"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}

// ListChildrenCategoriesHandler 列出子分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 05:09:47
func (s *categoryService) ListChildrenCategoriesHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id", "name"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	parentCategory, err := s.categoryDAO.GetByID(s.db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	if user.ID != parentCategory.UserID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	categories, pageInfo, err := s.categoryDAO.PaginateChildren(s.db, parentCategory, []string{"id", "name", "parent_id"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"categories": lo.Map(*categories, func(category model.Category, index int) map[string]interface{} {
				return category.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListChildrenArticlesHandler 列出子文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 01:38:12
func (s *categoryService) ListChildrenArticlesHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	parentCategory, err := s.categoryDAO.GetByID(s.db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	if user.ID != parentCategory.UserID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	articles, pageInfo, err := s.articleDAO.PaginateByCategoryID(s.db, parentCategory.ID, []string{"id", "title", "slug"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
