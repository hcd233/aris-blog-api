package handler

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

// ArticleVersionHandler 文章版本服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:54
type ArticleVersionHandler interface {
	HandleCreateArticleVersion(c *gin.Context)
	HandleGetArticleVersionInfo(c *gin.Context)
	HandleGetLatestArticleVersionInfo(c *gin.Context)
	HandleListArticleVersions(c *gin.Context)
}

type articleVersionHandler struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	articleDocDAO     *doc_dao.ArticleDocDAO
}

// NewArticleVersionHandler 创建文章版本服务
//
//	@return ArticleVersionHandler
//	@author centonhuang
//	@update 2024-12-08 16:59:54
func NewArticleVersionHandler() ArticleVersionHandler {
	return &articleVersionHandler{
		db:                database.GetDBInstance(),
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
	}
}

// HandleCreateArticleVersion 创建文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-17 12:44:17
func (h *articleVersionHandler) HandleCreateArticleVersion(c *gin.Context) {
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

	user, err := h.userDAO.GetByName(h.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	latestVersion, err := h.articleVersionDAO.GetLatestByArticleID(h.db, article.ID, []string{"version", "content"}, []string{})
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

	if err = h.articleVersionDAO.Create(h.db, articleVersion); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code: protocol.CodeCreateArticleVersionError,
		})
		return
	}

	if article.Status == model.ArticleStatusPublish {
		lo.Must0(h.articleDocDAO.UpdateDocument(document.TransformArticleToDocument(&model.Article{ID: article.ID, User: &model.User{}}, articleVersion)))

		if err := h.articleDAO.Update(h.db, article, map[string]interface{}{"published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
	}

	articleVersion = lo.Must1(h.articleVersionDAO.GetLatestByArticleID(h.db, article.ID, []string{"id", "created_at", "version"}, []string{}))

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articleVersion": articleVersion.GetBasicInfo(),
		},
	})
}

func (h *articleVersionHandler) HandleGetArticleVersionInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleVersionURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article version",
		})
		return
	}

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	articleVersion, err := h.articleVersionDAO.GetByArticleIDAndVersion(h.db, article.ID, uri.Version, []string{"id", "created_at", "version", "content", "summary"}, []string{})
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

func (h *articleVersionHandler) HandleGetLatestArticleVersionInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's article version",
		})
		return
	}

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	articleVersion, err := h.articleVersionDAO.GetLatestByArticleID(h.db, article.ID, []string{"id", "created_at", "version", "content", "summary"}, []string{})
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

// HandleListArticleVersions 列出文章版本
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-16 10:13:59
func (h *articleVersionHandler) HandleListArticleVersions(c *gin.Context) {
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

	user, err := h.userDAO.GetByName(h.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	versions, pageInfo, err := h.articleVersionDAO.PaginateByArticleID(h.db, article.ID, []string{"created_at", "version", "content"}, []string{}, param.Page, param.PageSize)
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
			"articleVersions": lo.Map(*versions, func(article model.ArticleVersion, _ int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}
