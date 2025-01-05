package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
)

// UserHandler 用户处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:20
type UserHandler interface {
	HandleGetCurUserInfo(c *gin.Context)
	HandleGetUserInfo(c *gin.Context)
	HandleUpdateInfo(c *gin.Context)
	HandleQueryUser(c *gin.Context)
}

type userHandler struct {
	svc service.UserService
}

// NewUserHandler 创建用户处理器
//
//	return UserHandler
//	author centonhuang
//	update 2024-12-08 16:59:38
func NewUserHandler() UserHandler {
	return &userHandler{
		svc: service.NewUserService(),
	}
}

func (h *userHandler) HandleGetCurUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	req := protocol.GetCurUserInfoRequest{
		CurUserID: userID,
	}

	rsp, err := h.svc.GetCurUserInfo(&req)

	util.SendHTTPResponse(c, rsp, err)
}

// GetUserInfoHandler 用户信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:30
func (h *userHandler) HandleGetUserInfo(c *gin.Context) {
	uri := c.MustGet("uri").(*protocol.UserURI)

	req := &protocol.GetUserInfoRequest{
		UserName: uri.UserName,
	}

	rsp, err := h.svc.GetUserInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// UpdateInfoHandler 更新用户信息
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:40
func (h *userHandler) HandleUpdateInfo(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	body := c.MustGet("body").(*protocol.UpdateUserBody)

	req := &protocol.UpdateUserInfoRequest{
		CurUserName:     userName,
		UserName:        uri.UserName,
		UpdatedUserName: body.UserName,
	}

	rsp, err := h.svc.UpdateUserInfo(req)

	util.SendHTTPResponse(c, rsp, err)
}

// QueryUserHandler 查询用户
//
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:50
func (h *userHandler) HandleQueryUser(c *gin.Context) {
	param := c.MustGet("param").(*protocol.QueryParam)

	req := &protocol.QueryUserRequest{
		QueryParam: param,
	}

	rsp, err := h.svc.QueryUser(req)

	util.SendHTTPResponse(c, rsp, err)
}
