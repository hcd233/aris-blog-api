package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TokenHandler 令牌处理器
//
//	author centonhuang
//	update 2025-01-05 21:00:00
type TokenHandler interface {
	HandleRefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*protocol.HumaHTTPResponse[*dto.RefreshTokenResponse], error)
}

type tokenHandler struct {
	svc service.TokenService
}

// NewTokenHandler 创建令牌处理器
//
//	return TokenHandler
//	author centonhuang
//	update 2025-01-05 21:00:00
func NewTokenHandler() TokenHandler {
	return &tokenHandler{
		svc: service.NewTokenService(),
	}
}

// HandleRefreshToken 刷新令牌
//
//	receiver h *tokenHandler
//	param ctx context.Context
//	param req *dto.RefreshTokenRequest
//	return *protocol.HumaHTTPResponse[*dto.RefreshTokenResponse]
//	return error
//	author centonhuang
//	update 2025-01-05 21:00:00
func (h *tokenHandler) HandleRefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*protocol.HumaHTTPResponse[*dto.RefreshTokenResponse], error) {
	return util.WrapHTTPResponse(h.svc.RefreshToken(ctx, req))
}
