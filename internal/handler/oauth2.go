package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// Oauth2Handler OAuth2处理器
type Oauth2Handler interface {
	HandleLogin(c *gin.Context)
	HandleCallback(c *gin.Context)
}

type githubOauth2Handler struct {
	svc service.Oauth2Service
}

// NewGithubOauth2Handler 创建Github OAuth2处理器
//
//	return Oauth2Handler
//	author centonhuang
//	update 2025-01-05 13:43:43
func NewGithubOauth2Handler() Oauth2Handler {
	return &githubOauth2Handler{
		svc: service.NewGithubOauth2Service(),
	}
}

// HandleLogin Github OAuth2登录
//
//	@Summary		Github OAuth2登录
//	@Description	Github OAuth2登录请求,返回重定向URL
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	protocol.HTTPResponse{data=protocol.LoginResponse,error=nil}
//	@Failure		400	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500	{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/github/login [get]
//	receiver h *githubOauth2Handler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 13:43:42
func (h *githubOauth2Handler) HandleLogin(c *gin.Context) {
	req := &protocol.LoginRequest{}

	rsp, err := h.svc.Login(c, req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleCallback Github OAuth2回调
//
//	@Summary		Github OAuth2回调
//	@Description	Github OAuth2回调请求,验证code和state
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"授权码"
//	@Param			state	query		string	true	"状态码"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.CallbackResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/oauth2/github/callback [get]
//	receiver h *githubOauth2Handler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-05 13:43:36
func (h *githubOauth2Handler) HandleCallback(c *gin.Context) {
	params := protocol.GithubCallbackParam{}
	if err := c.BindQuery(&params); err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return
	}

	req := &protocol.CallbackRequest{
		Code:  params.Code,
		State: params.State,
	}

	rsp, err := h.svc.Callback(c, req)

	util.SendHTTPResponse(c, rsp, err)
}
