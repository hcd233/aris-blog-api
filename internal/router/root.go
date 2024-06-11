// Package router provides the router implementation.
package router

import (
	"Aris-AI-go/internal/protocol"

	"github.com/gin-gonic/gin"
)

// GetRootMessage is the root message handler.
func GetRootMessage(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Message: "Welcome to Aris AI Go!",
		Status:  protocol.SUCCESS,
		Data:    nil,
	})
}
