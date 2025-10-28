package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type TokenHandlers struct{ svc service.TokenService }

func NewTokenHandlers() *TokenHandlers { return &TokenHandlers{svc: service.NewTokenService()} }

type refreshTokenInput struct{ humadto.RefreshTokenInput }

func (h *TokenHandlers) HandleRefreshToken(ctx context.Context, input *refreshTokenInput) (*humadto.Output[protocol.RefreshTokenResponse], error) {
	req := &protocol.RefreshTokenRequest{RefreshToken: input.Body.RefreshToken}
	rsp, err := h.svc.RefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.RefreshTokenResponse]{Body: *rsp}, nil
}
