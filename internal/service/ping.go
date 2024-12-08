package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

// PingService 健康检查服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type PingService interface {
	PingHandler(c *gin.Context)
}

type pingService struct{}

// NewPingService 创建健康检查服务
//
//	@return PingService
//	@author centonhuang
//	@update 2024-12-08 16:59:38
func NewPingService() PingService {
	return &pingService{}
}

func (s *pingService) PingHandler(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Welcome to Aris Blog API!",
	})
}
