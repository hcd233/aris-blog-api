// Package version 文章版本接口
//
//	@update 2024-10-16 10:10:30
package version

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListArticleVersionsHandler 列出文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-16 10:13:59
func ListArticleVersionsHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	param := c.MustGet("param").(*protocol.PageParam)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article versions",
		})
		return
	}

	db := database.GetDBInstance()

	userDAO, articleDAO, articleVersionDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetArticleVersionDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	versions, pageInfo, err := articleVersionDAO.PaginateByArticleID(db, article.ID, []string{"created_at", "version", "content"}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articleVersions": lo.Map(*versions, func(article model.ArticleVersion, index int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
