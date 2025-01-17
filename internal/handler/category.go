package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// CategoryHandler 分类服务
//
//	author centonhuang
//	update 2024-12-08 16:59:38
type CategoryHandler interface {
	HandleCreateCategory(c *gin.Context)
	HandleGetCategoryInfo(c *gin.Context)
	HandleUpdateCategoryInfo(c *gin.Context)
	HandleDeleteCategory(c *gin.Context)
	HandleGetRootCategories(c *gin.Context)
	HandleListChildrenCategories(c *gin.Context)
	HandleListChildrenArticles(c *gin.Context)
}

type categoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler 创建分类处理器
//
//	return CategoryHandler
//	author centonhuang
//	update 2024-12-08 16:5CategoryHandler
func NewCategoryHandler() CategoryHandler {
	return &categoryHandler{
		svc: service.NewCategoryService(),
	}
}

// CreateCategoryHandler 创建分类
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-09-28 07:03:28
func (h *categoryHandler) HandleCreateCategory(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.CreateCategoryBody)

	req := &protocol.CreateCategoryRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		Name:        body.Name,
		ParentID:    body.ParentID,
	}

	rsp, err := h.svc.CreateCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// GetCategoryInfoHandler 获取分类信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-01 04:58:27
func (h *categoryHandler) HandleGetCategoryInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	req := &protocol.GetCategoryInfoRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		CategoryID:  uri.CategoryID,
	}

	rsp, err := h.svc.GetCategoryInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetRootCategories 获取根分类信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-23 03:56:26
func (h *categoryHandler) HandleGetRootCategories(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)

	req := &protocol.GetRootCategoryRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
	}

	rsp, err := h.svc.GetRootCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateCategoryInfo 更新分类信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 03:45:55
func (h *categoryHandler) HandleUpdateCategoryInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	body := c.MustGet("body").(*protocol.UpdateCategoryBody)

	req := &protocol.UpdateCategoryRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		CategoryID:  uri.CategoryID,
		Name:        body.Name,
		ParentID:    body.ParentID,
	}

	rsp, err := h.svc.UpdateCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteCategory 删除分类
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 04:55:08
func (h *categoryHandler) HandleDeleteCategory(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	req := &protocol.DeleteCategoryRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		CategoryID:  uri.CategoryID,
	}

	rsp, err := h.svc.DeleteCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenCategories 列出子分类
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-01 05:09:47
func (h *categoryHandler) HandleListChildrenCategories(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenCategoriesRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		CategoryID:  uri.CategoryID,
		PageParam:   param,
	}

	rsp, err := h.svc.ListChildrenCategories(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenArticles 列出子文章
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 01:38:12
func (h *categoryHandler) HandleListChildrenArticles(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenArticlesRequest{
		CurUserName: userName,
		UserName:    uri.UserName,
		CategoryID:  uri.CategoryID,
		PageParam:   param,
	}

	rsp, err := h.svc.ListChildrenArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}
