// Package category 分类接口
//
//	@update 2024-09-23 11:42:18
package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListChildrenCategoriesHandler 列出子分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 05:09:47
func ListChildrenCategoriesHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	db := database.GetDBInstance()

	categoryDAO, userDAO := dao.GetCategoryDAO(), dao.GetUserDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id", "name"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	parentCategory, err := categoryDAO.GetByID(db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"})
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

	categories, pageInfo, err := categoryDAO.PaginateChildren(db, parentCategory, []string{"id", "name", "parent_id"}, param.Page, param.PageSize)
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
func ListChildrenArticlesHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's category",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, categoryDAO, articleDAO := dao.GetUserDAO(), dao.GetCategoryDAO(), dao.GetArticleDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	parentCategory, err := categoryDAO.GetByID(db, uri.CategoryID, []string{"id", "name", "parent_id", "user_id"})
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

	articles, pageInfo, err := articleDAO.PaginateByCategoryID(db, parentCategory.ID, []string{"id", "title", "slug"}, param.Page, param.PageSize)
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
