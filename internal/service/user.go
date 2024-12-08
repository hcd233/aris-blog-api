package service

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

// UserService 用户服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type UserService interface {
	GetMyInfoHandler(c *gin.Context)
	GetUserInfoHandler(c *gin.Context)
	UpdateInfoHandler(c *gin.Context)
	QueryUserHandler(c *gin.Context)
}

type userService struct {
	db            *gorm.DB
	userDAO       *dao.UserDAO
	tagDAO        *dao.TagDAO
	articleDAO    *dao.ArticleDAO
	userDocDAO    *doc_dao.UserDocDAO
	tagDocDAO     *doc_dao.TagDocDAO
	articleDocDAO *doc_dao.ArticleDocDAO
}

// NewUserService 创建用户服务
//
//	@return UserService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewUserService() UserService {
	return &userService{
		db:            database.GetDBInstance(),
		userDAO:       dao.GetUserDAO(),
		tagDAO:        dao.GetTagDAO(),
		articleDAO:    dao.GetArticleDAO(),
		userDocDAO:    doc_dao.GetUserDocDAO(),
		tagDocDAO:     doc_dao.GetTagDocDAO(),
		articleDocDAO: doc_dao.GetArticleDocDAO(),
	}
}

func (s *userService) GetMyInfoHandler(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := s.userDAO.GetByID(s.db, userID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
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
//	@update 2024-09-16 05:58:52
func (s *userService) GetUserInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
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
//	@update 2024-09-18 01:54:05
func (s *userService) UpdateInfoHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUserNotFoundError,
			Message: err.Error(),
		})
		return
	}

	lo.Must0(s.userDAO.Update(s.db, &model.User{ID: userID}, map[string]interface{}{
		"name": body.UserName,
	}))
	user = lo.Must1(s.userDAO.GetByID(s.db, userID, []string{"id", "name", "avatar"}, []string{}))

	var wg sync.WaitGroup
	var listTagErr, listArticleErr, updateTagDocErr, updateArticleDocErr, updateUserDocErr error
	var createdTags *[]model.Tag
	var createdArticles *[]model.Article

	wg.Add(2)

	go func() {
		defer wg.Done()
		createdTags, _, listTagErr = s.tagDAO.PaginateByUserID(s.db, userID, []string{"id"}, []string{}, 2, -1)
	}()
	go func() {
		defer wg.Done()
		createdArticles, _, listArticleErr = s.articleDAO.PaginateByUserID(s.db, userID, []string{"id"}, []string{}, 2, -1)
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
		updateTagDocErr = s.tagDocDAO.BatchUpdateDocuments(lo.Map(*createdTags, func(tag model.Tag, idx int) *document.TagDocument {
			return &document.TagDocument{ID: tag.ID, Creator: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateArticleDocErr = s.articleDocDAO.BatchUpdateDocuments(lo.Map(*createdArticles, func(article model.Article, idx int) *document.ArticleDocument {
			return &document.ArticleDocument{ID: article.ID, Author: user.Name}
		}))
	}()
	go func() {
		defer wg.Done()
		updateUserDocErr = s.userDocDAO.UpdateDocument(document.TransformUserToDocument(user))
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
//	@update 2024-09-16 05:58:52
func (s *userService) QueryUserHandler(c *gin.Context) {
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
