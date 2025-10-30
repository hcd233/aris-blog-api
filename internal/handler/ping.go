package handler

import (
	"context"
	"net/http"

	"github.com/hcd233/aris-blog-api/internal/protocol"
)

// PingHandler 健康检查处理器
//
//	author centonhuang
//	update 2025-01-04 15:52:48
type PingHandler interface {
	HandlePing(ctx context.Context, _ *struct{}) (*protocol.HumaResponse[*protocol.PingResponse], error)
}

type pingHandler struct{}

// NewPingHandler 创建健康检查处理器
//
//	return PingHandler
//	author centonhuang
//	update 2025-01-04 15:52:48
func NewPingHandler() PingHandler {
	return &pingHandler{}
}

// HandlePing 健康检查处理器
func (h *pingHandler) HandlePing(_ context.Context, _ *struct{}) (*protocol.HumaResponse[*protocol.PingResponse], error) {
	rsp := &protocol.PingResponse{
		Status: "ok",
	}

	return &protocol.HumaResponse[*protocol.PingResponse]{
		Status: http.StatusOK,
		Data:   rsp,
	}, nil
}
