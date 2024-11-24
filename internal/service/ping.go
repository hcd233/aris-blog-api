package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

type PingService interface {
	PingHandler(c *gin.Context)
}

type pingService struct {
}

func NewPingService() PingService {
	return &pingService{}
}

func (s *pingService) PingHandler(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Welcome to Aris Blog API!",
	})
}
