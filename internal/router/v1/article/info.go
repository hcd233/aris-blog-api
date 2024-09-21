// Package article 文章接口
//
//	@update 2024-09-21 05:37:21
package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// GetInfoHandler 文章信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleURI)

	article, err := model.QueryArticleBySlugAndUserName(uri.ArticleSlug, uri.UserName, []string{
		"id", "slug", "title", "status", "category",
		"created_at", "updated_at", "published_at",
		"likes", "views",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: article.GetDetailedInfo(),
	})
}

// UpdateArticleHandler 用户文章列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func UpdateArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	body := c.MustGet("body").(*protocol.UpdateArticleBody)

	updateFields := make(map[string]interface{})
	if body.Title != "" {
		updateFields["title"] = body.Title
	}
	if body.Slug != "" {
		updateFields["slug"] = body.Slug
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: "No any field to update",
		})
		return
	}

	article, err := model.QueryArticleBySlugAndUserName(uri.ArticleSlug, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	article, err = model.UpdateArticleInfoByID(article.ID, updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"article": article.GetBasicInfo(),
		},
	})
}

// DeleteArticleHandler 删除文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-22 04:32:37
func DeleteArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleURI)

	article, err := model.QueryArticleBySlugAndUserName(uri.ArticleSlug, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	err = article.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
