package version

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// CreateArticleVersionHandler 创建文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-17 12:44:17
func CreateArticleVersionHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ArticleURI)
	body := c.MustGet("body").(*protocol.CreateArticleVersionBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's article version",
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

	latestVersion, err := model.QueryLatestArticleVersionByArticleID(article.ID, []string{"version", "content"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	if latestVersion.Content == body.Content {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateArticleVersionError,
			Message: "The content of the new version is the same as the latest version",
		})
		return
	}

	articleVersion := model.ArticleVersion{
		Article: article,
		Content: body.Content,
		Version: latestVersion.Version + 1,
	}
	err = articleVersion.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeCreateArticleVersionError,
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: articleVersion.GetBasicInfo(),
	})
}
