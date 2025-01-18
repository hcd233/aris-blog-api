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
//	@Summary		创建分类
//	@Description	创建分类
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.CreateCategoryBody	true	"创建分类请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.CreateCategoryResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category [post]
//	param c *gin.Context
//	author centonhuang
//	update 2024-09-28 07:03:28
func (h *categoryHandler) HandleCreateCategory(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.CreateCategoryBody)

	req := &protocol.CreateCategoryRequest{
		UserID:   userID,
		Name:     body.Name,
		ParentID: body.ParentID,
	}

	rsp, err := h.svc.CreateCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetCategoryInfo 获取分类信息
//
//	@Summary		获取分类信息
//	@Description	根据分类ID获取分类详细信息
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			categoryID	path		uint	true	"分类ID"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.GetCategoryInfoResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/{categoryID} [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-01 04:58:27
func (h *categoryHandler) HandleGetCategoryInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	uri := c.MustGet("uri").(*protocol.CategoryURI)
	req := &protocol.GetCategoryInfoRequest{
		UserID:     userID,
		CategoryID: uri.CategoryID,
	}

	rsp, err := h.svc.GetCategoryInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetRootCategories 获取根分类信息
//
//	@Summary		获取根分类信息
//	@Description	获取根分类信息
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.GetRootCategoryResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/root [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-23 03:56:26
func (h *categoryHandler) HandleGetRootCategories(c *gin.Context) {
	userID := c.GetUint("userID")
	req := &protocol.GetRootCategoryRequest{
		UserID: userID,
	}

	rsp, err := h.svc.GetRootCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateCategoryInfo 更新分类信息
//
//	@Summary		更新分类信息
//	@Description	更新分类信息
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			path	path		protocol.CategoryURI	true	"分类ID"
//	@Param			body	body		protocol.UpdateCategoryBody	true	"更新分类请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.UpdateCategoryResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/{categoryID} [put]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 03:45:55
func (h *categoryHandler) HandleUpdateCategoryInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	body := c.MustGet("body").(*protocol.UpdateCategoryBody)

	req := &protocol.UpdateCategoryRequest{
		CategoryID: uri.CategoryID,
		Name:       body.Name,
		ParentID:   body.ParentID,
	}

	rsp, err := h.svc.UpdateCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteCategory 删除分类
//
//	@Summary		删除分类
//	@Description	删除分类
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			path	path		protocol.CategoryURI	true	"分类ID"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteCategoryResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/{categoryID} [delete]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 04:55:08
func (h *categoryHandler) HandleDeleteCategory(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.CategoryURI)

	req := &protocol.DeleteCategoryRequest{
		CategoryID: uri.CategoryID,
	}

	rsp, err := h.svc.DeleteCategory(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenCategories 列出子分类
//
//	@Summary		列出子分类
//	@Description	列出子分类
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			path	path		protocol.CategoryURI	true	"分类ID"
//	@Param			param	query		protocol.PageParam	    true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListChildrenCategoriesResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/{categoryID}/subCategories [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-01 05:09:47
func (h *categoryHandler) HandleListChildrenCategories(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenCategoriesRequest{
		UserID:     userID,
		CategoryID: uri.CategoryID,
		PageParam:  param,
	}

	rsp, err := h.svc.ListChildrenCategories(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListChildrenArticles 列出子文章
//
//	@Summary		列出子文章
//	@Description	列出子文章
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			path	path		protocol.CategoryURI	true	"分类ID"
//	@Param			param	query		protocol.PageParam	    true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListChildrenArticlesResponse,error=nil}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/category/{categoryID}/subArticles [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-10-02 01:38:12
func (h *categoryHandler) HandleListChildrenArticles(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.CategoryURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListChildrenArticlesRequest{
		UserID:     userID,
		CategoryID: uri.CategoryID,
		PageParam:  param,
	}

	rsp, err := h.svc.ListChildrenArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}
