package service

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

type ArticleService interface {
	CreateArticleHandler(c *gin.Context)
	GetArticleInfoHandler(c *gin.Context)
	UpdateArticleHandler(c *gin.Context)
	UpdateArticleStatusHandler(c *gin.Context)
	DeleteArticleHandler(c *gin.Context)
	ListArticlesHandler(c *gin.Context)
	ListUserArticlesHandler(c *gin.Context)
	QueryUserArticleHandler(c *gin.Context)
	QueryArticleHandler(c *gin.Context)
}

type articleService struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	tagDAO            *dao.TagDAO
	categoryDAO       *dao.CategoryDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
	articleDocDAO     *doc_dao.ArticleDocDAO
}

func NewArticleService() ArticleService {
	return &articleService{
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
func (s *articleService) CreateArticleHandler(c *gin.Context) {
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
		tag, err := s.tagDAO.GetBySlug(s.db, tagSlug, []string{"id"}, []string{})
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

	if err := s.articleDAO.Create(s.db, article); err != nil {
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
func (s *articleService) GetArticleInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.ArticleSlugURI)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{
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
func (s *articleService) UpdateArticleHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	articleDocDAO := doc_dao.GetArticleDocDAO()

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

	if err := s.articleDAO.Update(s.db, article, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateArticleError,
			Message: err.Error(),
		})
		return
	}

	article = lo.Must1(s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "title", "slug", "status"}, []string{}))
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
func (s *articleService) UpdateArticleStatusHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "status", "title", "slug", "user_id", "category_id"}, []string{"User", "Category", "Tags"})
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

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"content"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetArticleVersionError,
			Message: err.Error(),
		})
		return
	}

	if body.Status == model.ArticleStatusPublish {
		if err := s.articleDAO.Update(s.db, article, map[string]interface{}{"status": body.Status, "published_at": time.Now()}); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeUpdateArticleError,
				Message: err.Error(),
			})
			return
		}
		lo.Must0(articleDocDAO.AddDocument(document.TransformArticleToDocument(article, latestVersion)))
	} else if body.Status == model.ArticleStatusDraft {
		if err := s.articleDAO.Update(s.db, article, map[string]interface{}{"status": body.Status, "published_at": nil}); err != nil {
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
func (s *articleService) DeleteArticleHandler(c *gin.Context) {
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

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, uri.ArticleSlug, user.ID, []string{"id", "slug"}, []string{})
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
		deleteArticleErr = s.articleDAO.Delete(s.db, article)
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
func (s *articleService) ListArticlesHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	articles, pageInfo, err := s.articleDAO.PaginateByPublished(
		s.db,
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
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
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
func (s *articleService) ListUserArticlesHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryUserError,
			Message: err.Error(),
		})
		return
	}

	articles, pageInfo, err := s.articleDAO.PaginateByUserID(s.db, user.ID, []string{"id", "title", "slug"}, []string{}, param.Page, param.PageSize)
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
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
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
func (s *articleService) QueryUserArticleHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.QueryParam)

	_, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	articles, queryInfo, err := s.articleDocDAO.QueryDocument(param.Query, append(param.Filter, "author="+uri.UserName), param.Page, param.PageSize)
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
func (s *articleService) QueryArticleHandler(c *gin.Context) {
	params := c.MustGet("param").(*protocol.QueryParam)

	articles, queryInfo, err := s.articleDocDAO.QueryDocument(params.Query, params.Filter, params.Page, params.PageSize)
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
