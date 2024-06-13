// Package router provides the router implementation.
package router

import (
	"github.com/hcd233/Aris-AI-go/internal/protocol"

	"github.com/gin-gonic/gin"
)

// RootHandler is the root message handler.
func rootHandler(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Message: "Welcome to Aris AI Go!",
		Status:  protocol.SUCCESS,
		Data:    nil,
	})
}
