package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type Oauth2Handlers struct{ svc service.Oauth2Service }

func NewGithubOauth2Handlers() *Oauth2Handlers {
	return &Oauth2Handlers{svc: service.NewGithubOauth2Service()}
}
func NewQQOauth2Handlers() *Oauth2Handlers { return &Oauth2Handlers{svc: service.NewQQOauth2Service()} }
func NewGoogleOauth2Handlers() *Oauth2Handlers {
	return &Oauth2Handlers{svc: service.NewGoogleOauth2Service()}
}

type (
	providerPathInput  struct{ humadto.ProviderPath }
	oauthCallbackQuery struct {
		Code  string `query:"code"`
		State string `query:"state"`
	}
)

func (h *Oauth2Handlers) HandleLogin(ctx context.Context, _ *providerPathInput) (*humadto.Output[protocol.LoginResponse], error) {
	req := &protocol.LoginRequest{}
	rsp, err := h.svc.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.LoginResponse]{Body: *rsp}, nil
}

func (h *Oauth2Handlers) HandleCallback(ctx context.Context, input *oauthCallbackQuery) (*humadto.Output[protocol.CallbackResponse], error) {
	req := &protocol.CallbackRequest{Code: input.Code, State: input.State}
	rsp, err := h.svc.Callback(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.CallbackResponse]{Body: *rsp}, nil
}
