package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ListArticlesHandler 列出文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func ListArticlesHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	articleDAO := dao.GetArticleDAO()

	articles, pageInfo, err := articleDAO.PaginateByPublished(
		db,
		[]string{"id", "title", "slug", "status", "published_at", "views", "likes", "user_id"},
		[]string{"User", "Comments", "Tags"},
		param.Page, param.PageSize,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
				return article.GetDetailedInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserArticlesHandler 用户文章列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func ListUserArticlesHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	userDAO, articleDAO := dao.GetUserDAO(), dao.GetArticleDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	articles, pageInfo, err := articleDAO.PaginateByUserID(db, user.ID, []string{"id", "title", "slug"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
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
