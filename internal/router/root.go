// Package router provides the router implementation.
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/protocol"
)

// RootHandler is the root message handler.
func handleRoot(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Message: "Welcome to Aris AI Go!",
		Status:  protocol.SUCCESS,
		Data:    nil,
	})
}
