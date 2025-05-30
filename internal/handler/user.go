package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// UserHandler 用户处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:20
type UserHandler interface {
	HandleGetCurUserInfo(c *gin.Context)
	HandleGetUserInfo(c *gin.Context)
	HandleUpdateInfo(c *gin.Context)
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

// HandleGetCurUserInfo 获取当前用户信息
//
//	@Summary		获取当前用户信息
//	@Description	获取当前用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.GetCurUserInfoResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user/current [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:30
func (h *userHandler) HandleGetCurUserInfo(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)

	req := &protocol.GetCurUserInfoRequest{
		UserID: userID,
	}

	rsp, err := h.svc.GetCurUserInfo(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// GetUserInfoHandler 用户信息
//
//	@Summary		获取用户信息
//	@Description	获取用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			path	path		protocol.UserURI	true	"用户名"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.GetUserInfoResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user/{userID} [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:30
func (h *userHandler) HandleGetUserInfo(c *gin.Context) {
	uri := c.MustGet(constant.CtxKeyURI).(*protocol.UserURI)

	req := &protocol.GetUserInfoRequest{
		UserID: uri.UserID,
	}

	rsp, err := h.svc.GetUserInfo(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// UpdateInfoHandler 更新用户信息
//
//	@Summary		更新用户信息
//	@Description	更新用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		protocol.UpdateUserBody	true	"更新用户信息请求"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.UpdateUserInfoResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/user [patch]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:40
func (h *userHandler) HandleUpdateInfo(c *gin.Context) {
	userID := c.GetUint(constant.CtxKeyUserID)
	body := c.MustGet(constant.CtxKeyBody).(*protocol.UpdateUserBody)

	req := &protocol.UpdateUserInfoRequest{
		UserID:          userID,
		UpdatedUserName: body.UserName,
	}

	rsp, err := h.svc.UpdateUserInfo(c, req)

	util.SendHTTPResponse(c, rsp, err)
}
