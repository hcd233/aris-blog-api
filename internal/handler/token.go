package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// TokenHandler 令牌处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:10
type TokenHandler interface {
	HandleRefreshToken(c *gin.Context)
}

type tokenHandler struct {
	svc service.TokenService
}

// NewTokenHandler 创建令牌处理器
//
//	return TokenHandler
//	author centonhuang
//	update 2025-01-04 15:56:04
func NewTokenHandler() TokenHandler {
	return &tokenHandler{
		svc: service.NewTokenService(),
	}
}

// HandleRefreshToken 刷新令牌
//
//	receiver s *tokenHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:56:10
func (h *tokenHandler) HandleRefreshToken(c *gin.Context) {
	body := c.MustGet("body").(*protocol.RefreshTokenBody)

	req := &protocol.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	}

	rsp, err := h.svc.RefreshToken(req)

	util.SendHTTPResponse(c, rsp, err)
}
