// Package user 用户接口
//
//	@update 2024-09-16 05:29:08
package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/resource/search"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

// GetUserInfoHandler 用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func GetUserInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)

	db := database.GetDBInstance()

	userDAO := dao.GetUserDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: fmt.Sprintf("User `%s` not found", uri.UserName),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: user.GetDetailedInfo(),
	})
}

// UpdateInfoHandler 更新用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-18 01:54:05
func UpdateInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.UpdateUserBody)
	userID, userName := c.GetUint("userID"), c.GetString("userName")

	db := database.GetDBInstance()
	searchEngine := search.GetSearchEngine()

	userDAO, tagDAO, articleDAO := dao.GetUserDAO(), dao.GetTagDAO(), dao.GetArticleDAO()
	userDocDAO, tagDocDAO, articleDocDAO := docdao.GetUserDocDAO(), docdao.GetTagDocDAO(), docdao.GetArticleDocDAO()

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's info",
		})
		return
	}

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(userDAO.Update(db, &model.User{ID: userID}, map[string]interface{}{
		"name": body.UserName,
	}))
	user = lo.Must1(userDAO.GetByID(db, userID, []string{"id", "name", "avatar"}))

	createdTags := lo.Must1(tagDAO.ListByUserID(db, userID, []string{"id"}, -1, -1))
	createdArticles := lo.Must1(articleDAO.ListByUserID(db, userID, []string{"id"}, -1, -1))

	lo.Must0(tagDocDAO.BatchUpdateDocuments(searchEngine, lo.Map(*createdTags, func(tag model.Tag, idx int) *document.TagDocument {
		return &document.TagDocument{ID: tag.ID, Creator: user.Name}
	})))
	lo.Must0(articleDocDAO.BatchUpdateDocuments(searchEngine, lo.Map(*createdArticles, func(article model.Article, idx int) *document.ArticleDocument {
		return &document.ArticleDocument{ID: article.ID, Author: user.Name}
	})))
	lo.Must0(userDocDAO.UpdateDocument(searchEngine, document.TransformUserToDocument(user)))

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Update user info successfully",
	})
}
