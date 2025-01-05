package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
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
	HandleListUserTags(c *gin.Context)
	HandleQueryTag(c *gin.Context)
	HandleQueryUserTag(c *gin.Context)
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

func (h *tagHandler) HandleCreateTag(c *gin.Context) {
	userID := c.GetUint("userID")
	body := c.MustGet("body").(*protocol.CreateTagBody)

	req := protocol.CreateTagRequest{
		CurUserID:   userID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.CreateTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// GetTagInfoHandler 获取标签信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:52:48
func (h *tagHandler) HandleGetTagInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.TagURI)

	req := protocol.GetTagInfoRequest{
		TagSlug: uri.TagSlug,
	}

	rsp, err := h.svc.GetTagInfo(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUpdateTag 更新标签
//
//	receiver s *tagHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:16
func (h *tagHandler) HandleUpdateTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)
	body := c.MustGet("body").(*protocol.UpdateTagBody)

	req := protocol.UpdateTagRequest{
		CurUserID:   userID,
		TagSlug:     uri.TagSlug,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
	}

	rsp, err := h.svc.UpdateTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteTag 删除标签
//
//	receiver s *tagHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:24
func (h *tagHandler) HandleDeleteTag(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.TagURI)

	req := protocol.DeleteTagRequest{
		CurUserID: userID,
		TagName:   uri.TagSlug,
	}

	rsp, err := h.svc.DeleteTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListTags 标签列表
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:31
func (h *tagHandler) HandleListTags(c *gin.Context) {
	param := c.MustGet("param").(*protocol.PageParam)

	req := protocol.ListTagsRequest{
		PageParam: param,
	}

	rsp, err := h.svc.ListTags(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserTags 列出用户标签
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:38
func (h *tagHandler) HandleListUserTags(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	req := protocol.ListUserTagsRequest{
		UserName:  uri.UserName,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserTags(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleQueryTag 搜索标签
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:45
func (h *tagHandler) HandleQueryTag(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	req := protocol.QueryTagRequest{
		QueryParam: param,
	}

	rsp, err := h.svc.QueryTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleQueryUserTag 搜索用户标签
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:55:52
func (h *tagHandler) HandleQueryUserTag(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)
	params := c.MustGet("param").(*protocol.QueryParam)

	req := protocol.QueryUserTagRequest{
		UserName:   uri.UserName,
		QueryParam: params,
	}

	rsp, err := h.svc.QueryUserTag(&req)

	util.SendHTTPResponse(c, rsp, err)
}
