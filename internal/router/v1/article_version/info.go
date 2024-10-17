package version

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
)

// GetArticleVersionInfoHandler 获取文章版本信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetArticleVersionInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleVersionURI)

	db := database.GetDBInstance()

	userDAO, articleDAO, articleVersionDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetArticleVersionDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
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

	articleVersion, err := articleVersionDAO.GetByArticleIDAndVersion(db, article.ID, uri.Version, []string{"id", "created_at", "version", "content"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: articleVersion.GetDetailedInfo(),
	})
}
