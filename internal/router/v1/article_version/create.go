package version

import (
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
	"github.com/samber/lo"
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

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id", "status"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"version", "content"})
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

	articleVersion := &model.ArticleVersion{
		Article: article,
		Content: body.Content,
		Version: latestVersion.Version + 1,
	}

	if err = articleVersionDAO.Create(db, articleVersion); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeCreateArticleVersionError,
		})
		return
	}

	if article.Status == model.ArticleStatusPublish {
		searchEngine := search.GetSearchEngine()
		articleDocDAO := docdao.GetArticleDocDAO()

		lo.Must0(articleDocDAO.UpdateDocument(searchEngine, document.TransformArticleToDocument(&model.Article{ID: article.ID, User: &model.User{}}, articleVersion)))

		if err := articleDAO.Update(db, article, map[string]interface{}{"published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: articleVersion.GetBasicInfo(),
	})
}
