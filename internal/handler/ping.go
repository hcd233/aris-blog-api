package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

// PingHandler 健康检查处理器
//
//	@author centonhuang
//	@update 2025-01-04 15:52:48
type PingHandler interface {
	HandlePing(c *gin.Context)
}

type pingHandler struct{}

// NewPingHandler 创建健康检查处理器
//
//	@return PingHandler
//	@author centonhuang
//	@update 2025-01-04 15:52:48
func NewPingHandler() PingHandler {
	return &pingHandler{}
}

func (h *pingHandler) HandlePing(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Welcome to Aris Blog API!",
	})
}
