package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/oauth2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// Oauth2Handler OAuth2处理器接口
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type Oauth2Handler interface {
	HandleLogin(ctx context.Context, req *dto.LoginRequest) (*protocol.HTTPResponse[*dto.LoginResponse], error)
	HandleCallback(ctx context.Context, req *dto.CallbackRequest) (*protocol.HTTPResponse[*dto.CallbackResponse], error)
}

type oauth2Handler struct{}

// NewOauth2Handler 创建OAuth2处理器
//
//	return Oauth2Handler
//	author centonhuang
//	update 2025-01-05 21:00:00
func NewOauth2Handler() Oauth2Handler {
	return &oauth2Handler{}
}

// HandleLogin OAuth2登录
//	@receiver h *oauth2Handler 
//	@param ctx context.Context 
//	@param req *dto.LoginRequest 
//	@return *protocol.HTTPResponse[*dto.LoginResponse] 
//	@return error 
//	@author centonhuang 
//	@update 2025-11-02 04:16:14 
func (h *oauth2Handler) HandleLogin(ctx context.Context, req *dto.LoginRequest) (*protocol.HTTPResponse[*dto.LoginResponse], error) {
	svc := h.getService(oauth2.ProviderType(req.Provider))
	return util.WrapHTTPResponse(svc.Login(ctx, req))
}

// HandleCallback OAuth2回调
//	@receiver h *oauth2Handler 
//	@param ctx context.Context 
//	@param req *dto.CallbackRequest 
//	@return *protocol.HTTPResponse[*dto.CallbackResponse] 
//	@return error 
//	@author centonhuang 
//	@update 2025-11-02 04:16:22 
func (h *oauth2Handler) HandleCallback(ctx context.Context, req *dto.CallbackRequest) (*protocol.HTTPResponse[*dto.CallbackResponse], error) {
	svc := h.getService(oauth2.ProviderType(req.Provider))
	return util.WrapHTTPResponse(svc.Callback(ctx, req))
}

// getService 根据provider获取对应的service
//
//	receiver h *oauth2Handler
//	param provider string
//	return service.Oauth2Service
//	author centonhuang
//	update 2025-01-05 21:00:00
func (h *oauth2Handler) getService(provider oauth2.ProviderType) service.Oauth2Service {
	switch provider {
	case oauth2.ProviderTypeGithub:
		return service.NewGithubOauth2Service()
	case oauth2.ProviderTypeGoogle:
		return service.NewGoogleOauth2Service()
	default:
		return service.NewGithubOauth2Service() // 默认返回 github
	}
}
