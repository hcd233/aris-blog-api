package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
//
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-16 04:07:30
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://example.com"},    // 设置允许的域
		AllowMethods:     []string{"GET", "POST", "PUT"},     // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type"}, // 允许的Header
		ExposeHeaders:    []string{"Content-Length"},         // 公开的Header
		AllowCredentials: true,                               // 是否允许cookie
		MaxAge:           12 * time.Hour,
	})
}
