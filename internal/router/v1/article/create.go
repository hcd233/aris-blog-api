package article

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// CreateArticleHandler 创建文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-21 09:58:14
func CreateArticleHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateArticleBody)

	db := database.GetDBInstance()

	tagDAO, articleDAO := dao.GetTagDAO(), dao.GetArticleDAO()

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
		tag, err := tagDAO.GetBySlug(db, tagSlug, []string{"id"}, []string{})
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

	if err := articleDAO.Create(db, article); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateArticleError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: article.GetBasicInfo(),
	})
}
