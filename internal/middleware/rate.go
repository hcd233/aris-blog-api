package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterMiddleware 限频中间件
//
//	@param period time.Duration
//	@param limit int64
//	@param key string
//	@param errCode protocol.ResponseCode
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-10-22 05:01:51
func RateLimiterMiddleware(period time.Duration, limit int64, key string, errorCode protocol.ResponseCode) gin.HandlerFunc {
	// 创建限频规则
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	// 使用内存存储限频数据
	store := memory.NewStore()

	// 创建限频实例
	instance := limiter.New(store, rate)

	// 创建中间件
	middleware := mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		c.JSON(http.StatusTooManyRequests, protocol.Response{
			Code: errorCode,
		})
	}))

	return func(c *gin.Context) {
		// 获取限频 key
		value := c.Param(key)
		if key == "" {
			value = c.ClientIP() // 如果没有指定的参数，则使用 IP 地址作为 key
		}

		// 设置限频 key
		c.Set("limiter", fmt.Sprintf("%s:%s", key, value))

		// 应用限频中间件
		middleware(c)

		// 继续处理请求
		c.Next()
	}
}
