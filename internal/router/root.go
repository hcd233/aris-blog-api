// Package router provides the router implementation.
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

// RootHandler 根路由
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-09-16 02:07:30
func RootHandler(c *gin.Context) {
	c.JSON(200, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Welcome to Aris Blog API!",
	})
}
