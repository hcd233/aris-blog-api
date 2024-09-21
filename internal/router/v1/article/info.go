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
