package handler

import (
	"net/http"
	"sync"
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

// ArticleHandler 文章服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type ArticleHandler interface {
	HandleCreateArticle(c *gin.Context)
	HandleGetArticleInfo(c *gin.Context)
	HandleUpdateArticle(c *gin.Context)
	HandleUpdateArticleStatus(c *gin.Context)
	HandleDeleteArticle(c *gin.Context)
	HandleListArticles(c *gin.Context)
	HandleListUserArticles(c *gin.Context)
	HandleQueryUserArticle(c *gin.Context)
	HandleQueryArticle(c *gin.Context)
}

type articleHandler struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	tagDAO            *dao.TagDAO
	categoryDAO       *dao.CategoryDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	articleDocDAO     *doc_dao.ArticleDocDAO
}

// NewArticleHandler 创建文章服务
//
//	@return ArticleService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewArticleHandler() ArticleHandler {
	return &articleHandler{
		db:          database.GetDBInstance(),
		userDAO:     dao.GetUserDAO(),
		tagDAO:      dao.GetTagDAO(),
		categoryDAO: dao.GetCategoryDAO(),
		articleDAO:  dao.GetArticleDAO(),
	}
}

// CreateArticleHandler 创建文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 09:58:14
func (h *articleHandler) HandleCreateArticle(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateArticleBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's article",
		})
		return
	}

	if body.Slug == "" {
		body.Slug = body.Title
	}

	tags := []model.Tag{}
	tagChan, errChan := make(chan *model.Tag, len(body.Tags)), make(chan error, len(body.Tags))

	var wg sync.WaitGroup
	wg.Add(len(body.Tags))

	getTagFunc := func(tagSlug string) {
		defer wg.Done()
		tag, err := h.tagDAO.GetBySlug(h.db, tagSlug, []string{"id"}, []string{})
		if err != nil {
			errChan <- err
			return
		}
		tagChan <- tag
	}

	for _, tagSlug := range body.Tags {
		go getTagFunc(tagSlug)
	}

	wg.Wait()
	close(tagChan)
	close(errChan)

	if len(errChan) > 0 {
		err := <-errChan
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	for tag := range tagChan {
		tags = append(tags, *tag)
	}

	article := &model.Article{
		UserID:     userID,
		Status:     model.ArticleStatusDraft,
		Title:      body.Title,
		Slug:       body.Slug,
		Tags:       tags,
		CategoryID: body.CategoryID,
		Comments:   []model.Comment{},
		Versions:   []model.ArticleVersion{},
	}

	if err := h.articleDAO.Create(h.db, article); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"article": article.GetBasicInfo(),
		},
	})
}

// GetArticleInfoHandler 文章信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 05:58:52
func (h *articleHandler) HandleGetArticleInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{
		"id", "slug", "title", "status", "user_id",
		"created_at", "updated_at", "published_at",
		"likes", "views",
	}, []string{"User", "Comments", "Tags"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"article": article.GetDetailedInfo(),
		},
	})
}

// UpdateArticleHandler 用户文章列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func (h *articleHandler) HandleUpdateArticle(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleBody)
	userName := c.GetString("userName")

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's article",
		})
		return
	}

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

	user, err := h.userDAO.GetByName(h.db, userName, []string{"id"}, []string{})
	articleDocDAO := doc_dao.GetArticleDocDAO()

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

	if err := h.articleDAO.Update(h.db, article, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: err.Error(),
		})
		return
	}

	article = lo.Must1(h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "title", "slug", "status"}, []string{}))
	if article.Status == model.ArticleStatusPublish {
		article.User = &model.User{}
		lo.Must0(articleDocDAO.UpdateDocument(document.TransformArticleToDocument(article, &model.ArticleVersion{})))
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
func (h *articleHandler) HandleUpdateArticleStatus(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)
	body := c.MustGet("body").(*protocol.UpdateArticleStatusBody)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's article status",
		})
		return
	}

	articleDocDAO := doc_dao.GetArticleDocDAO()

	user, err := h.userDAO.GetByName(h.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "status", "title", "slug", "user_id", "category_id"}, []string{"User", "Category", "Tags"})
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

	latestVersion, err := h.articleVersionDAO.GetLatestByArticleID(h.db, article.ID, []string{"content"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	if body.Status == model.ArticleStatusPublish {
		if err := h.articleDAO.Update(h.db, article, map[string]interface{}{"status": body.Status, "published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
		lo.Must0(articleDocDAO.AddDocument(document.TransformArticleToDocument(article, latestVersion)))
	} else if body.Status == model.ArticleStatusDraft {
		if err := h.articleDAO.Update(h.db, article, map[string]interface{}{"status": body.Status, "published_at": nil}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
		lo.Must0(articleDocDAO.DeleteDocument(article.ID))
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
func (h *articleHandler) HandleDeleteArticle(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's article",
		})
		return
	}

	articleDocDAO := doc_dao.GetArticleDocDAO()

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleDAO.GetBySlugAndUserID(h.db, uri.ArticleSlug, user.ID, []string{"id", "slug"}, []string{})
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
		deleteArticleErr = h.articleDAO.Delete(h.db, article)
	}()

	go func() {
		defer wg.Done()
		deleteDocErr = articleDocDAO.DeleteDocument(article.ID)
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

// ListArticlesHandler 列出文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 08:59:40
func (h *articleHandler) HandleListArticles(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	articles, pageInfo, err := h.articleDAO.PaginateByPublished(
		h.db,
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
			"articles": lo.Map(*articles, func(article model.Article, _ int) map[string]interface{} {
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
func (h *articleHandler) HandleListUserArticles(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	articles, pageInfo, err := h.articleDAO.PaginateByUserID(h.db, user.ID, []string{"id", "title", "slug"}, []string{}, param.Page, param.PageSize)
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
			"articles": lo.Map(*articles, func(article model.Article, _ int) map[string]interface{} {
				return article.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// QueryUserArticleHandler 查询用户文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:13:21
func (h *articleHandler) HandleQueryUserArticle(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.QueryParam)

	_, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	articles, queryInfo, err := h.articleDocDAO.QueryDocument(param.Query, append(param.Filter, "author="+uri.UserName), param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles":  articles,
			"queryInfo": queryInfo,
		},
	})
}

// QueryArticleHandler 查询文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:13:17
func (h *articleHandler) HandleQueryArticle(c *gin.Context) {
	params := c.MustGet("param").(*protocol.QueryParam)

	articles, queryInfo, err := h.articleDocDAO.QueryDocument(params.Query, params.Filter, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles":  articles,
			"queryInfo": queryInfo,
		},
	})
}
