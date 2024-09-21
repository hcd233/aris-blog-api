package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// CreateArticleHandler 创建文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 09:58:14
func CreateArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateArticleBody)

	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)

	if uri.UserName != userName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code: protocol.CodeNotPermissionError,
		})
		return
	}

	if body.Slug == "" {
		body.Slug = body.Title
	}

	article := model.Article{
		UserID: userID,
		Status: model.ArticleStatusDraft,
		Title:  body.Title,
		Slug:   body.Slug,
	}

	err := article.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeCreateArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: article.GetBasicInfo(),
	})
}
