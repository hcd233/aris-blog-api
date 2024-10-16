// Package category 分类接口
//
//	@update 2024-09-23 11:42:18
package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListRootCategoriesHandler 列出根分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-23 11:45:08
func ListRootCategoriesHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := model.QueryUserByName(uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	rootCategory, err := model.QueryRootCategoryByUserID(user.ID, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetCategoryError,
			Message: err.Error(),
		})
		return
	}

	categories, err := model.QueryChildrenCategoriesByUserID(rootCategory.ID, []string{"id", "name", "parent_id"}, param.Limit, param.Offset)
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
		},
	})
}

// ListChildrenCategoriesHandler 列出子分类
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 05:09:47
func ListChildrenCategoriesHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	categories, err := model.QueryChildrenCategoriesByUserID(uri.CategoryID, []string{"id", "name", "parent_id"}, param.Limit, param.Offset)
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
		},
	})
}

// ListChildrenArticlesHandler 列出子文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-02 01:38:12
func ListChildrenArticlesHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	articles, err := model.QueryChildrenArticlesByCategoryID(uri.CategoryID, []string{"id", "title", "slug"}, param.Limit, param.Offset)
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
		},
	})
}
