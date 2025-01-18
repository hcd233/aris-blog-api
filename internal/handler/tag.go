package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TagHandler 标签处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type TagHandler interface {
	HandleCreateTag(c *gin.Context)
	HandleGetTagInfo(c *gin.Context)
	HandleUpdateTag(c *gin.Context)
	HandleDeleteTag(c *gin.Context)
	HandleListTags(c *gin.Context)
}

type tagHandler struct {
	svc service.TagService
}

// NewTagHandler 创建标签处理器
//
//	return TagHandler
//	author centonhuang
//	update 2025-01-04 15:52:48
func NewTagHandler() TagHandler {
	return &tagHandler{
		svc: service.NewTagService(),
	}
}

// HandleCreateTag 创建标签
//
//	@Summary		创建标签
//	@Description	创建标签
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			body	body		protocol.CreateTagBody	true	"创建标签请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.CreateTagResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tag [post]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:52:48
func (h *tagHandler) HandleCreateTag(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.CreateTagBody)

	req := protocol.CreateTagRequest{
		UserID:      userID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.CreateTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetTagInfo 获取标签信息
//
//	@Summary		获取标签信息
//	@Description	根据标签slug获取标签详细信息
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			path	path	protocol.TagURI	true	"标签ID"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.GetTagInfoResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tag/{tagSlug} [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:52:48
func (h *tagHandler) HandleGetTagInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	req := &protocol.GetTagInfoRequest{
		TagID: uri.TagID,
	}

	rsp, err := h.svc.GetTagInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateTag 更新标签
//
//	@Summary		更新标签
//	@Description	更新标签
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			path	path	protocol.TagURI         true	"标签ID"
//	@Param			body	body	protocol.UpdateTagBody	true	"更新标签请求体"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.UpdateTagResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tag/{tagSlug} [put]
//	receiver s *tagHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:16
func (h *tagHandler) HandleUpdateTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	req := &protocol.UpdateTagRequest{
		UserID:      userID,
		TagID:       uri.TagID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.UpdateTag(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteTag 删除标签
//
//	@Summary		删除标签
//	@Description	删除标签
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			tagSlug	path		string	true	"标签slug"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteTagResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tag/{tagSlug} [delete]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:24
func (h *tagHandler) HandleDeleteTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)

	req := &protocol.DeleteTagRequest{
		UserID: userID,
		TagID:  uri.TagID,
	}

	rsp, err := h.svc.DeleteTag(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListTags 列出标签
//
//	@Summary		列出标签
//	@Description	获取标签列表
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			param	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListTagsResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tags [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:31
func (h *tagHandler) HandleListTags(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListTagsRequest{
		PageParam: param,
	}

	rsp, err := h.svc.ListTags(req)

	util.SendHTTPResponse(c, rsp, err)
}
