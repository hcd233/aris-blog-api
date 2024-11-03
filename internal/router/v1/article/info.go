// Package article 文章接口
//
//	@update 2024-09-21 05:37:21
package article

import (
	"net/http"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

// GetArticleInfoHandler 文章信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetArticleInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	db := database.GetDBInstance()

	userDAO, articleDAO := dao.GetUserDAO(), dao.GetArticleDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{
		"id", "slug", "title", "status", "category_id",
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
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleBody)
	userName := c.MustGet("userName").(string)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's article",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, articleDAO := dao.GetUserDAO(), dao.GetArticleDAO()

	updateFields := make(map[string]interface{})

	if body.Title != "" {
		updateFields["title"] = body.Title
	}
	if body.Slug != "" {
		updateFields["slug"] = body.Slug
	}

	if body.CategoryID != 0 {
		updateFields["category_id"] = body.CategoryID
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: "No any field to update",
		})
		return
	}

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	articleDocDAO := docdao.GetArticleDocDAO()

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

	if err := articleDAO.Update(db, article, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: err.Error(),
		})
		return
	}

	article = lo.Must1(articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id", "title", "slug", "status"}))
	if article.Status == model.ArticleStatusPublish {
		article.User = &model.User{}
		lo.Must0(articleDocDAO.UpdateDocument(searchEngine, document.TransformArticleToDocument(article, &model.ArticleVersion{})))
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"article": article.GetBasicInfo(),
		},
	})
}

// UpdateArticleStatusHandler 更新文章状态
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-17 09:28:54
func UpdateArticleStatusHandler(c *gin.Context) {
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleStatusBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's article status",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, articleDAO, articleVersionDAO := dao.GetUserDAO(), dao.GetArticleDAO(), dao.GetArticleVersionDAO()

	articleDocDAO := docdao.GetArticleDocDAO()

	user, err := userDAO.GetByName(db, userName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetAllBySlugAndUserID(db, uri.ArticleSlug, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if article.Status == body.Status {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: "The status is the same as the current status",
		})
		return
	}

	latestVersion, err := articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"content"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	if body.Status == model.ArticleStatusPublish {
		if err := articleDAO.Update(db, article, map[string]interface{}{"status": body.Status, "published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
		lo.Must0(articleDocDAO.AddDocument(searchEngine, document.TransformArticleToDocument(article, latestVersion)))
	} else if body.Status == model.ArticleStatusDraft {
		if err := articleDAO.Update(db, article, map[string]interface{}{"status": body.Status, "published_at": nil}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
		lo.Must0(articleDocDAO.DeleteDocument(searchEngine, article.ID))
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
	userName := c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's article",
		})
		return
	}

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, articleDAO := dao.GetUserDAO(), dao.GetArticleDAO()

	articleDocDAO := docdao.GetArticleDocDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := articleDAO.GetBySlugAndUserID(db, uri.ArticleSlug, user.ID, []string{"id", "slug"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	var wg sync.WaitGroup
	var deleteArticleErr, deleteDocErr error
	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteArticleErr = articleDAO.Delete(db, article)
	}()

	go func() {
		defer wg.Done()
		deleteDocErr = articleDocDAO.DeleteDocument(searchEngine, article.ID)
	}()

	if deleteArticleErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: deleteArticleErr.Error(),
		})
		return
	}

	if deleteDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteArticleError,
			Message: deleteDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
