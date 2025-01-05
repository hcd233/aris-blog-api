package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/util"
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

// HandlePing 处理健康检查请求
//
//	@receiver h *pingHandler
//	@param c *gin.Context
//	@author centonhuang
//	@update 2025-01-04 20:47:48
func (h *pingHandler) HandlePing(c *gin.Context) {
	rsp := protocol.PingResponse{
		Status: "ok",
	}

	util.SendHTTPResponse(c, rsp, nil)
}
