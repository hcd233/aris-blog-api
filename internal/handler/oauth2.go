package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
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
//	@return Oauth2Handler
//	@author centonhuang
//	@update 2025-01-05 13:43:43
func NewGithubOauth2Handler() Oauth2Handler {
	return &githubOauth2Handler{
		svc: service.NewGithubOauth2Service(),
	}
}

// HandleLogin 处理Github OAuth2登录
//
//	@receiver h *githubOauth2Handler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 13:43:42
func (h *githubOauth2Handler) HandleLogin(c *gin.Context) {
	req := &protocol.LoginRequest{}

	rsp, err := h.svc.Login(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleCallback 处理Github OAuth2回调
//
//	@receiver h *githubOauth2Handler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-05 13:43:36
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

	rsp, err := h.svc.Callback(req)

	util.SendHTTPResponse(c, rsp, err)
}
