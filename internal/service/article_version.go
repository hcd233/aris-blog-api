package service

import (
	"errors"
	"net/http"
	"time"

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

type ArticleVersionService interface {
	CreateArticleVersionHandler(c *gin.Context)
	GetArticleVersionInfoHandler(c *gin.Context)
	GetLatestArticleVersionInfoHandler(c *gin.Context)
	ListArticleVersionsHandler(c *gin.Context)
}

type articleVersionService struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	articleDocDAO     *doc_dao.ArticleDocDAO
}

func NewArticleVersionService() ArticleVersionService {
	return &articleVersionService{
		db:                database.GetDBInstance(),
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
	}
}

// CreateArticleVersionHandler 创建文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-17 12:44:17
func (s *articleVersionService) CreateArticleVersionHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.CreateArticleVersionBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's article version",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"version", "content"}, []string{})
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

	if err = s.articleVersionDAO.Create(s.db, articleVersion); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeCreateArticleVersionError,
		})
		return
	}

	if article.Status == model.ArticleStatusPublish {
		lo.Must0(s.articleDocDAO.UpdateDocument(document.TransformArticleToDocument(&model.Article{ID: article.ID, User: &model.User{}}, articleVersion)))

		if err := s.articleDAO.Update(s.db, article, map[string]interface{}{"published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
	}

	articleVersion = lo.Must1(s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "created_at", "version"}, []string{}))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articleVersion": articleVersion.GetBasicInfo(),
		},
	})
}

// GetArticleVersionInfoHandler 获取文章版本信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func (s *articleVersionService) GetArticleVersionInfoHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleVersionURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article version",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	articleVersion, err := s.articleVersionDAO.GetByArticleIDAndVersion(s.db, article.ID, uri.Version, []string{"id", "created_at", "version", "content", "summary"}, []string{})
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
			"articleVersion": articleVersion.GetDetailedInfo(),
		},
	})
}

func (s *articleVersionService) GetLatestArticleVersionInfoHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article version",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	articleVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"id", "created_at", "version", "content", "summary"}, []string{})
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
			"articleVersion": articleVersion.GetDetailedInfo(),
		},
	})
}

// ListArticleVersionsHandler 列出文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-16 10:13:59
func (s *articleVersionService) ListArticleVersionsHandler(c *gin.Context) {
	userName := c.GetString("userName")
	param := c.MustGet("param").(*protocol.PageParam)
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article versions",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	versions, pageInfo, err := s.articleVersionDAO.PaginateByArticleID(s.db, article.ID, []string{"created_at", "version", "content"}, []string{}, param.Page, param.PageSize)
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
