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

// TagService 标签服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type TagService interface {
	CreateTagHandler(c *gin.Context)
	GetTagInfoHandler(c *gin.Context)
	UpdateTagHandler(c *gin.Context)
	DeleteTagHandler(c *gin.Context)
	ListTagsHandler(c *gin.Context)
	ListUserTagsHandler(c *gin.Context)
	QueryTagHandler(c *gin.Context)
	QueryUserTagHandler(c *gin.Context)
}

type tagService struct {
	db        *gorm.DB
	userDAO   *dao.UserDAO
	tagDAO    *dao.TagDAO
	tagDocDAO *doc_dao.TagDocDAO
}

// NewTagService 创建标签服务
//
//	@return TagService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewTagService() TagService {
	return &tagService{
		db:        database.GetDBInstance(),
		userDAO:   dao.GetUserDAO(),
		tagDAO:    dao.GetTagDAO(),
		tagDocDAO: doc_dao.GetTagDocDAO(),
	}
}

func (s *tagService) CreateTagHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	body := c.MustGet("body").(*protocol.CreateTagBody)

	tag := &model.Tag{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		User:        &model.User{ID: userID, Name: userName},
	}

	var wg sync.WaitGroup
	var createTagErr, createDocErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		createTagErr = s.tagDAO.Create(s.db, tag)
	}()

	go func() {
		defer wg.Done()
		createDocErr = s.tagDocDAO.AddDocument(document.TransformTagToDocument(tag))
	}()

	wg.Wait()

	if createTagErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: createTagErr.Error(),
		})
		return
	}

	if createDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateTagError,
			Message: createDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tag": tag.GetBasicInfo(),
		},
	})
}

// GetTagInfoHandler 获取标签信息
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-01 04:58:01
func (s *tagService) GetTagInfoHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	tag, err := s.tagDAO.GetBySlug(s.db, uri.TagSlug, []string{"id", "name", "slug", "description", "created_by"}, []string{"User"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tag": tag.GetDetailedInfo(),
		},
	})
}

func (s *tagService) UpdateTagHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	tag, err := s.tagDAO.GetBySlug(s.db, uri.TagSlug, []string{"id", "created_by"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	if tag.CreatedBy != userID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to update other user's tag",
		})
		return
	}

	updateFields := make(map[string]interface{})
	if body.Name != "" {
		updateFields["name"] = body.Name
	}
	if body.Slug != "" {
		updateFields["slug"] = body.Slug
	}
	if body.Description != "" {
		updateFields["description"] = body.Description
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: "No fields to update",
		})
		return
	}

	if err := s.tagDAO.Update(s.db, tag, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	tag = lo.Must1(s.tagDAO.GetBySlug(s.db, uri.TagSlug, []string{"id", "name", "slug", "description", "created_by"}, []string{"User"}))

	// 同步到搜索引擎
	err = s.tagDocDAO.UpdateDocument(document.TransformTagToDocument(tag))
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tag": tag.GetDetailedInfo(),
		},
	})
}

// DeleteTagHandler 删除标签
func (s *tagService) DeleteTagHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)

	tag, err := s.tagDAO.GetBySlug(s.db, uri.TagSlug, []string{"id", "name", "slug", "created_by"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetTagError,
			Message: err.Error(),
		})
		return
	}

	if tag.CreatedBy != userID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's tag",
		})
		return
	}

	var wg sync.WaitGroup
	var deleteTagErr, deleteDocErr error
	wg.Add(2)

	go func() {
		defer wg.Done()
		deleteTagErr = s.tagDAO.Delete(s.db, tag)
	}()

	go func() {
		defer wg.Done()
		deleteDocErr = s.tagDocDAO.DeleteDocument(tag.ID)
	}()

	wg.Wait()

	if deleteTagErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: deleteTagErr.Error(),
		})
		return
	}

	if deleteDocErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteTagError,
			Message: deleteDocErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}

// ListTagsHandler 标签列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-22 02:41:01
func (s *tagService) ListTagsHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	tags, pageInfo, err := s.tagDAO.Paginate(s.db, []string{"id", "slug", "name"}, []string{}, param.Page, param.PageSize)
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
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserTagsHandler 列出用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-22 02:41:01
func (s *tagService) ListUserTagsHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	db := database.GetDBInstance()

	userDAO, tagDAO := dao.GetUserDAO(), dao.GetTagDAO()

	user, err := userDAO.GetByName(db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags, pageInfo, err := tagDAO.PaginateByUserID(db, user.ID, []string{"id", "slug", "name"}, []string{}, param.Page, param.PageSize)
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
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// QueryTagHandler 搜索标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:35:21
func (s *tagService) QueryTagHandler(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	tags, queryInfo, err := s.tagDocDAO.QueryDocument(param.Query, param.Filter, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags":      tags,
			"queryInfo": queryInfo,
		},
	})
}

// QueryUserTagHandler 搜索用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-10-23 12:35:31
func (s *tagService) QueryUserTagHandler(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	params := c.MustGet("param").(*protocol.QueryParam)

	_, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags, queryInfo, err := s.tagDocDAO.QueryDocument(params.Query, append(params.Filter, "creator="+uri.UserName), params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeQueryTagError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags":      tags,
			"queryInfo": queryInfo,
		},
	})
}
