// Package user 用户接口
//
//	@update 2024-09-16 05:29:08
package user

import (
	"fmt"
	"net/http"
	"sync"

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

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
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

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
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
	user = lo.Must1(userDAO.GetByID(db, userID, []string{"id", "name", "avatar"}, []string{}))

	var wg sync.WaitGroup
	var listTagErr, listArticleErr, updateTagDocErr, updateArticleDocErr, updateUserDocErr error
	var createdTags *[]model.Tag
	var createdArticles *[]model.Article

	wg.Add(2)

	go func() {
		defer wg.Done()
		createdTags, _, listTagErr = tagDAO.PaginateByUserID(db, userID, []string{"id"}, []string{}, 2, -1)
	}()
	go func() {
		defer wg.Done()
		createdArticles, _, listArticleErr = articleDAO.PaginateByUserID(db, userID, []string{"id"}, []string{}, 2, -1)
	}()

	wg.Wait()

	if listTagErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: listTagErr.Error(),
		})
		return
	}
	if listArticleErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: listArticleErr.Error(),
		})
		return
	}

	wg.Add(3)

	go func() {
		defer wg.Done()
		updateTagDocErr = tagDocDAO.BatchUpdateDocuments(searchEngine, lo.Map(*createdTags, func(tag model.Tag, idx int) *document.TagDocument {
			return &document.TagDocument{ID: tag.ID, Creator: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateArticleDocErr = articleDocDAO.BatchUpdateDocuments(searchEngine, lo.Map(*createdArticles, func(article model.Article, idx int) *document.ArticleDocument {
			return &document.ArticleDocument{ID: article.ID, Author: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateUserDocErr = userDocDAO.UpdateDocument(searchEngine, document.TransformUserToDocument(user))
	}()

	wg.Wait()

	if updateTagDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: updateTagDocErr.Error(),
		})
		return
	}

	if updateArticleDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: updateArticleDocErr.Error(),
		})
		return
	}

	if updateUserDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateUserError,
			Message: updateUserDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Update user info successfully",
	})
}
