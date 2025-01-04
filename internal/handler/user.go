package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	doc_dao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:56:20
type UserHandler interface {
	HandleGetMyInfo(c *gin.Context)
	HandleGetUserInfo(c *gin.Context)
	HandleUpdateInfo(c *gin.Context)
	HandleQueryUser(c *gin.Context)
}

type userHandler struct {
	db            *gorm.DB
	userDAO       *dao.UserDAO
	tagDAO        *dao.TagDAO
	articleDAO    *dao.ArticleDAO
	userDocDAO    *doc_dao.UserDocDAO
	tagDocDAO     *doc_dao.TagDocDAO
	articleDocDAO *doc_dao.ArticleDocDAO
}

// NewUserHandler 创建用户处理器
//
//	@return UserHandler
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewUserHandler() UserHandler {
	return &userHandler{
		db:            database.GetDBInstance(),
		userDAO:       dao.GetUserDAO(),
		tagDAO:        dao.GetTagDAO(),
		articleDAO:    dao.GetArticleDAO(),
		userDocDAO:    doc_dao.GetUserDocDAO(),
		tagDocDAO:     doc_dao.GetTagDocDAO(),
		articleDocDAO: doc_dao.GetArticleDocDAO(),
	}
}

func (h *userHandler) HandleGetMyInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userDAO.GetByID(h.db, userID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"user": user.GetDetailedInfo(),
		},
	})
}

// GetUserInfoHandler 用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:56:30
func (h *userHandler) HandleGetUserInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"user": user.GetDetailedInfo(),
		},
	})
}

// UpdateInfoHandler 更新用户信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:56:40
func (h *userHandler) HandleUpdateInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.UpdateUserBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's info",
		})
		return
	}

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(h.userDAO.Update(h.db, &model.User{ID: userID}, map[string]interface{}{
		"name": body.UserName,
	}))
	user = lo.Must1(h.userDAO.GetByID(h.db, userID, []string{"id", "name", "avatar"}, []string{}))

	var wg sync.WaitGroup
	var listTagErr, listArticleErr, updateTagDocErr, updateArticleDocErr, updateUserDocErr error
	var createdTags *[]model.Tag
	var createdArticles *[]model.Article

	wg.Add(2)

	go func() {
		defer wg.Done()
		createdTags, _, listTagErr = h.tagDAO.PaginateByUserID(h.db, userID, []string{"id"}, []string{}, 2, -1)
	}()
	go func() {
		defer wg.Done()
		createdArticles, _, listArticleErr = h.articleDAO.PaginateByUserID(h.db, userID, []string{"id"}, []string{}, 2, -1)
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
		updateTagDocErr = h.tagDocDAO.BatchUpdateDocuments(lo.Map(*createdTags, func(tag model.Tag, _ int) *document.TagDocument {
			return &document.TagDocument{ID: tag.ID, Creator: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateArticleDocErr = h.articleDocDAO.BatchUpdateDocuments(lo.Map(*createdArticles, func(article model.Article, _ int) *document.ArticleDocument {
			return &document.ArticleDocument{ID: article.ID, Author: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateUserDocErr = h.userDocDAO.UpdateDocument(document.TransformUserToDocument(user))
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

// QueryUserHandler 查询用户
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:56:50
func (h *userHandler) HandleQueryUser(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	docDAO := doc_dao.GetUserDocDAO()

	users, queryInfo, err := docDAO.QueryDocument(param.Query, param.Filter, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"users":     users,
			"queryInfo": queryInfo,
		},
	})
}
