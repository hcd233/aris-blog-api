package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TagHandler 标签处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type TagHandler interface {
	HandleCreateTag(c *fiber.Ctx) error
	HandleGetTagInfo(c *fiber.Ctx) error
	HandleUpdateTag(c *fiber.Ctx) error
	HandleDeleteTag(c *fiber.Ctx) error
	HandleListTags(c *fiber.Ctx) error
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
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-04 15:52:48
func (h *tagHandler) HandleCreateTag(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	body := c.Locals(constant.CtxKeyBody).(*protocol.CreateTagBody)

	req := protocol.CreateTagRequest{
		UserID:      userID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.CreateTag(c.Context(), &req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
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
//	@Router			/v1/tag/{tagID} [get]
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-04 15:52:48
func (h *tagHandler) HandleGetTagInfo(c *fiber.Ctx) error {
	uri := c.Locals(constant.CtxKeyURI).(*protocol.TagURI)

	req := &protocol.GetTagInfoRequest{
		TagID: uri.TagID,
	}

	rsp, err := h.svc.GetTagInfo(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
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
//	@Router			/v1/tag/{tagID} [patch]
//	receiver s *tagHandler
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-04 15:55:16
func (h *tagHandler) HandleUpdateTag(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.TagURI)
	body := c.Locals(constant.CtxKeyBody).(*protocol.UpdateTagBody)

	req := &protocol.UpdateTagRequest{
		UserID:      userID,
		TagID:       uri.TagID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.UpdateTag(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}

// HandleDeleteTag 删除标签
//
//	@Summary		删除标签
//	@Description	删除标签
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			path	path	protocol.TagURI	true	"标签ID"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteTagResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/tag/{tagID} [delete]
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-04 15:55:24
func (h *tagHandler) HandleDeleteTag(c *fiber.Ctx) error {
	userID := c.Locals(constant.CtxKeyUserID).(uint)
	uri := c.Locals(constant.CtxKeyURI).(*protocol.TagURI)

	req := &protocol.DeleteTagRequest{
		UserID: userID,
		TagID:  uri.TagID,
	}

	rsp, err := h.svc.DeleteTag(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
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
//	@Router			/v1/tag/list [get]
//	param c *fiber.Ctx error
//	author centonhuang
//	update 2025-01-04 15:55:31
func (h *tagHandler) HandleListTags(c *fiber.Ctx) error {
	param := c.Locals(constant.CtxKeyParam).(*protocol.PageParam)

	req := &protocol.ListTagsRequest{
		PageParam: param,
	}

	rsp, err := h.svc.ListTags(c.Context(), req)

	util.SendHTTPResponse(c, rsp, err)
	return nil
}
