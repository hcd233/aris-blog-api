package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
)

// ValidateURIMiddleware 验证URI中间件
//
//	@param uri interface{}
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-21 07:47:53
func ValidateURIMiddleware(uri interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindUri(uri); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeURIError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("uri", uri)
		c.Next()
	}
}

// ValidateParamMiddleware 验证参数中间件
//
//	@param param interface{}
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-21 07:48:40
func ValidateParamMiddleware(param interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(param); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeURIError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("param", param)
		c.Next()
	}
}

// ValidateBodyMiddleware 验证请求体中间件
//
//	@param body interface{}
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-21 08:48:25
func ValidateBodyMiddleware(body interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, protocol.Response{
				Code:    protocol.CodeBodyError,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("body", body)
		c.Next()
	}
}
