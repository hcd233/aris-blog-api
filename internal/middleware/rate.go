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

// RateLimiterConfig 限频配置
type RateLimiterConfig struct {
	Period    time.Duration         // 限频周期
	Limit     int64                 // 限制次数
	Key       string                // 用于生成限频key
	ErrorCode protocol.ResponseCode // 限频错误码
}

// RateLimiterMiddleware 限频中间件
//
//	@param config RateLimiterConfig
//	@return gin.HandlerFunc
//	@author AliceAI
//	@update 2023-05-22 15:30:00
func RateLimiterMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	// 创建限频规则
	rate := limiter.Rate{
		Period: config.Period,
		Limit:  config.Limit,
	}

	// 使用内存存储限频数据
	store := memory.NewStore()

	// 创建限频实例
	instance := limiter.New(store, rate)

	// 创建中间件
	middleware := mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		c.JSON(http.StatusTooManyRequests, protocol.Response{
			Code: config.ErrorCode,
		})
		
	}))

	return func(c *gin.Context) {
		// 获取限频 key
		key := c.Param(config.Key)
		if key == "" {
			key = c.ClientIP() // 如果没有指定的参数，则使用 IP 地址作为 key
		}

		// 设置限频 key
		c.Set("limiter", fmt.Sprintf("%s:%s", config.Key, key))

		// 应用限频中间件
		middleware(c)

		// 继续处理请求
		c.Next()
	}
}
