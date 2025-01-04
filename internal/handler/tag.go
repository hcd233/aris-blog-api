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

// TagHandler 标签处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:52:48
type TagHandler interface {
	HandleCreateTag(c *gin.Context)
	HandleGetTagInfo(c *gin.Context)
	HandleUpdateTag(c *gin.Context)
	HandleDeleteTag(c *gin.Context)
	HandleListTags(c *gin.Context)
	HandleListUserTags(c *gin.Context)
	HandleQueryTag(c *gin.Context)
	HandleQueryUserTag(c *gin.Context)
}

type tagHandler struct {
	db        *gorm.DB
	userDAO   *dao.UserDAO
	tagDAO    *dao.TagDAO
	tagDocDAO *doc_dao.TagDocDAO
}

// NewTagHandler 创建标签处理器
//
//	@return TagHandler
//	@author centonhuang
//	@update 2025-01-04 15:52:48
func NewTagHandler() TagHandler {
	return &tagHandler{
		db:        database.GetDBInstance(),
		userDAO:   dao.GetUserDAO(),
		tagDAO:    dao.GetTagDAO(),
		tagDocDAO: doc_dao.GetTagDocDAO(),
	}
}

func (h *tagHandler) HandleCreateTag(c *gin.Context) {
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
		createTagErr = h.tagDAO.Create(h.db, tag)
	}()

	go func() {
		defer wg.Done()
		createDocErr = h.tagDocDAO.AddDocument(document.TransformTagToDocument(tag))
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
//	@update 2025-01-04 15:52:48
func (h *tagHandler) HandleGetTagInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	tag, err := h.tagDAO.GetBySlug(h.db, uri.TagSlug, []string{"id", "name", "slug", "description", "created_by"}, []string{"User"})
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

// HandleUpdateTag 更新标签
//
//	@receiver s *tagHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:16
func (h *tagHandler) HandleUpdateTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	tag, err := h.tagDAO.GetBySlug(h.db, uri.TagSlug, []string{"id", "created_by"}, []string{})
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

	if err := h.tagDAO.Update(h.db, tag, updateFields); err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUpdateTagError,
			Message: err.Error(),
		})
		return
	}

	tag = lo.Must1(h.tagDAO.GetBySlug(h.db, uri.TagSlug, []string{"id", "name", "slug", "description", "created_by"}, []string{"User"}))

	// 同步到搜索引擎
	err = h.tagDocDAO.UpdateDocument(document.TransformTagToDocument(tag))
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

// HandleDeleteTag 删除标签
//
//	@receiver s *tagHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:24
func (h *tagHandler) HandleDeleteTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)

	tag, err := h.tagDAO.GetBySlug(h.db, uri.TagSlug, []string{"id", "name", "slug", "created_by"}, []string{})
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
		deleteTagErr = h.tagDAO.Delete(h.db, tag)
	}()

	go func() {
		defer wg.Done()
		deleteDocErr = h.tagDocDAO.DeleteDocument(tag.ID)
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

// HandleListTags 标签列表
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:31
func (h *tagHandler) HandleListTags(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	tags, pageInfo, err := h.tagDAO.Paginate(h.db, []string{"id", "slug", "name"}, []string{}, param.Page, param.PageSize)
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
			"tags": lo.Map(*tags, func(tag model.Tag, _ int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// HandleListUserTags 列出用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:38
func (h *tagHandler) HandleListUserTags(c *gin.Context) {
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
			"tags": lo.Map(*tags, func(tag model.Tag, _ int) map[string]interface{} {
				return tag.GetBasicInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// HandleQueryTag 搜索标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:45
func (h *tagHandler) HandleQueryTag(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	tags, queryInfo, err := h.tagDocDAO.QueryDocument(param.Query, param.Filter, param.Page, param.PageSize)
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

// HandleQueryUserTag 搜索用户标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 15:55:52
func (h *tagHandler) HandleQueryUserTag(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	params := c.MustGet("param").(*protocol.QueryParam)

	_, err := h.userDAO.GetByName(h.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	tags, queryInfo, err := h.tagDocDAO.QueryDocument(params.Query, append(params.Filter, "creator="+uri.UserName), params.Page, params.PageSize)
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
